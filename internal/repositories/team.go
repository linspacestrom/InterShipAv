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

type TeamRepo interface {
	Create(ctx context.Context, team domain.Team) (domain.Team, error)
	GetByName(ctx context.Context, name string) (domain.Team, error)
}

type TeamRepository struct {
	pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{pool: pool}
}

func (t *TeamRepository) Create(ctx context.Context, team domain.Team) (domain.Team, error) {
	var createdTeam domain.Team

	tx := transaction.GetQuerier(ctx, t.pool)

	query := `INSERT INTO team (team_name) VALUES ($1) RETURNING team_name`
	row := tx.QueryRow(ctx, query, team.Name)

	if err := row.Scan(&createdTeam.Name); err != nil {
		return createdTeam, err
	}

	return createdTeam, nil
}

func (t *TeamRepository) GetByName(ctx context.Context, name string) (domain.Team, error) {
	var team domain.Team

	tx := transaction.GetQuerier(ctx, t.pool)

	query := `SELECT * FROM team WHERE team.team_name = $1`
	row := tx.QueryRow(ctx, query, name)

	if err := row.Scan(&team.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return team, validateError.TeamNotFound
		}
		return team, err
	}

	return team, nil
}
