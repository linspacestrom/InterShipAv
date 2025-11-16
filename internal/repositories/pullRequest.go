package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/linspacestrom/InterShipAv/internal/domain"
	"github.com/linspacestrom/InterShipAv/internal/transaction"
	"github.com/linspacestrom/InterShipAv/internal/validateError"
)

type PrRepo interface {
	Create(ctx context.Context, pr domain.PullRequestCreate) (domain.PullRequestRead, error)
	Merge(ctx context.Context, prId string) (domain.PRMergeRead, error)
	GetById(ctx context.Context, id string) (domain.PullRequestRead, error)
	AssignReviewers(ctx context.Context, prId string, userIds []string) ([]string, error)
	GetPullRequestIdsByReviewerId(ctx context.Context, userId string) ([]string, error)
	GetPullRequestsByIds(ctx context.Context, prIds []string) ([]domain.PullRequestReviewRead, error)
	GetReviewersById(ctx context.Context, id string) ([]string, error)
}

type PullRequestRepository struct {
	pool *pgxpool.Pool
}

func NewPullRequestRepository(pool *pgxpool.Pool) *PullRequestRepository {
	return &PullRequestRepository{pool: pool}
}

func (r *PullRequestRepository) Create(ctx context.Context, createPR domain.PullRequestCreate) (domain.PullRequestRead, error) {
	var pr domain.PullRequestRead

	tx := transaction.GetQuerier(ctx, r.pool)

	row := tx.QueryRow(ctx, `INSERT INTO pull_request (pull_request_id, pull_request_name, author_id, status) 
			VALUES ($1, $2, $3, $4) RETURNING pull_request_id, pull_request_name, author_id, status`, createPR.Id, createPR.Name, createPR.AuthorId, domain.StatusOpen)

	if err := row.Scan(&pr.Id, &pr.Name, &pr.AuthorId, &pr.Status); err != nil {
		return pr, err
	}

	return pr, nil

}

func (r *PullRequestRepository) Merge(ctx context.Context, prId string) (domain.PRMergeRead, error) {
	tx := transaction.GetQuerier(ctx, r.pool)

	row := tx.QueryRow(ctx, `
		UPDATE pull_request 
		SET 
			status = $1,
			merged_at = COALESCE(merged_at, NOW())
		WHERE pull_request_id = $2
		RETURNING pull_request_id, pull_request_name, author_id, status, merged_at;
	`, domain.StatusMerged, prId)

	var prMerged domain.PRMergeRead

	if err := row.Scan(&prMerged.Id, &prMerged.Name, &prMerged.AuthorId, &prMerged.Status, &prMerged.MergedAt); err != nil {
		return prMerged, err
	}

	return prMerged, nil
}

func (r *PullRequestRepository) GetById(ctx context.Context, id string) (domain.PullRequestRead, error) {
	var pr domain.PullRequestRead

	q := transaction.GetQuerier(ctx, r.pool)
	row := q.QueryRow(ctx, `SELECT pull_request_id, pull_request_name, author_id, status FROM pull_request WHERE pull_request_id = $1`, id)

	if err := row.Scan(&pr.Id, &pr.Name, &pr.AuthorId, &pr.Status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pr, validateError.ErrPrNotExist
		}
		return pr, err
	}
	return pr, nil
}

func (r *PullRequestRepository) AssignReviewers(ctx context.Context, prId string, userIds []string) ([]string, error) {
	reviewerIds := make([]string, 0, 2)

	tx := transaction.GetQuerier(ctx, r.pool)

	for _, userId := range userIds {
		var reviewerId string
		row := tx.QueryRow(ctx, `INSERT INTO pr_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2) RETURNING reviewer_id`, prId, userId)
		if err := row.Scan(&reviewerId); err != nil {
			return reviewerIds, err
		}
		reviewerIds = append(reviewerIds, reviewerId)
	}

	return reviewerIds, nil
}

func (r *PullRequestRepository) GetPullRequestIdsByReviewerId(ctx context.Context, userId string) ([]string, error) {
	tx := transaction.GetQuerier(ctx, r.pool)

	rows, err := tx.Query(ctx, `SELECT pull_request_id FROM pr_reviewers WHERE reviewer_id = $1`, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var prIds []string
	for rows.Next() {
		var prId string
		if err := rows.Scan(&prId); err != nil {
			return nil, err
		}
		prIds = append(prIds, prId)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return prIds, nil
}

func (r *PullRequestRepository) GetPullRequestsByIds(ctx context.Context, prIds []string) ([]domain.PullRequestReviewRead, error) {
	tx := transaction.GetQuerier(ctx, r.pool)

	rows, err := tx.Query(ctx, `SELECT pull_request_id, pull_request_name, author_id, status FROM pull_request WHERE pull_request_id = ANY($1) AND status = $2`, prIds, domain.StatusOpen)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prReviews := make([]domain.PullRequestReviewRead, 0, len(prIds))

	for rows.Next() {
		var pr domain.PullRequestReviewRead
		if err := rows.Scan(&pr.Id, &pr.Name, &pr.AuthorId, &pr.Status); err != nil {
			return nil, err
		}
		prReviews = append(prReviews, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return prReviews, nil
}

func (r *PullRequestRepository) GetReviewersById(ctx context.Context, id string) ([]string, error) {
	reviewerIds := make([]string, 0, 2)

	tx := transaction.GetQuerier(ctx, r.pool)

	rows, err := tx.Query(ctx, `SELECT reviewer_id FROM pr_reviewers WHERE pull_request_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var reviewerId string
		if err := rows.Scan(&reviewerId); err != nil {
			return nil, err
		}
		reviewerIds = append(reviewerIds, reviewerId)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviewerIds, nil

}
