package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/linspacestrom/InterShipAv/internal/domain"
	"github.com/linspacestrom/InterShipAv/internal/transaction"
	"github.com/linspacestrom/InterShipAv/internal/validateError"
)

type UserRepo interface {
	AddUsersToTeam(ctx context.Context, users []domain.TeamMember, teamName string) error
	GetUserByTeamName(ctx context.Context, name string) ([]domain.TeamMember, error)
	GetById(ctx context.Context, id string) (domain.User, error)
	SetActiveById(ctx context.Context, id string, isActive bool) (domain.User, error)
	GetReviewers(ctx context.Context, name string, excludeUserId string) ([]string, error)
}

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) AddUsersToTeam(ctx context.Context, users []domain.TeamMember, teamName string) error {
	tx := transaction.GetQuerier(ctx, r.pool)

	valueStrings := make([]string, 0, len(users))
	valueArgs := make([]any, 0, len(users)*4)

	arg := 1

	log.Printf("\nпользователи создались\n")

	for _, u := range users {
		valueStrings = append(
			valueStrings,
			fmt.Sprintf("($%d, $%d, $%d, $%d)", arg, arg+1, arg+2, arg+3),
		)
		valueArgs = append(valueArgs, u.ID, u.Username, teamName, u.IsActive)
		arg += 4
	}

	query := fmt.Sprintf(`
		INSERT INTO "user" (id, username, team_name, is_active)
		VALUES %s
		ON CONFLICT (id) DO UPDATE SET
			username = EXCLUDED.username,
			team_name = EXCLUDED.team_name,
			is_active = EXCLUDED.is_active
	`, strings.Join(valueStrings, ","))

	_, err := tx.Exec(ctx, query, valueArgs...)
	return err
}

func (r *UserRepository) GetUserByTeamName(ctx context.Context, name string) ([]domain.TeamMember, error) {
	tx := transaction.GetQuerier(ctx, r.pool)

	rows, err := tx.Query(ctx, `SELECT id, username, is_active FROM "user" WHERE team_name = $1`, name)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []domain.TeamMember
	for rows.Next() {
		var u domain.TeamMember
		if err := rows.Scan(&u.ID, &u.Username, &u.IsActive); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetById(ctx context.Context, id string) (domain.User, error) {
	var user domain.User

	q := transaction.GetQuerier(ctx, r.pool)
	row := q.QueryRow(ctx, `SELECT id, username, team_name, is_active FROM "user" WHERE id = $1`, id)

	if err := row.Scan(&user.Id, &user.Username, &user.TeamName, &user.IsActive); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, validateError.UserNotFound
		}
		return user, err
	}
	return user, nil
}

func (r *UserRepository) SetActiveById(ctx context.Context, id string, isActive bool) (domain.User, error) {
	var user domain.User

	tx := transaction.GetQuerier(ctx, r.pool)

	row := tx.QueryRow(ctx, `UPDATE "user" SET is_active = $1 WHERE id = $2 RETURNING id, username, team_name, is_active`, isActive, id)

	if err := row.Scan(&user.Id, &user.Username, &user.TeamName, &user.IsActive); err != nil {
		return user, err
	}

	return user, nil

}

func (r *UserRepository) GetReviewers(ctx context.Context, name string, excludeUserId string) ([]string, error) {
	userIds := make([]string, 0, 2)

	tx := transaction.GetQuerier(ctx, r.pool)

	rows, err := tx.Query(ctx, `SELECT id FROM "user" WHERE is_active = true and id != $1 and team_name = $2 ORDER BY RANDOM() LIMIT 2`, excludeUserId, name)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIds = append(userIds, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userIds, nil
}
