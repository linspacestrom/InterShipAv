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
	GetById(ctx context.Context, id string) (domain.PullRequestRead, error)
	AssignReviewers(ctx context.Context, prId string, userIds []string) ([]string, error)
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

	row := tx.QueryRow(ctx, `INSERT INTO pull_request (pull_request_id, pull_request_name, author_id, status) VALUES ($1, $2, $3, $4) RETURNING pull_request_id, pull_request_name, author_id, status`, createPR.Id, createPR.Name, createPR.AuthorId, domain.StatusOpen)

	if err := row.Scan(&pr.Id, &pr.Name, &pr.AuthorId, &pr.Status); err != nil {
		return pr, err
	}

	return pr, nil

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
			return reviewerIds, nil
		}
		reviewerIds = append(reviewerIds, reviewerId)
	}

	return reviewerIds, nil
}
