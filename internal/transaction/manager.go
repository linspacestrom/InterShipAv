package transaction

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

type Manager struct {
	pool *pgxpool.Pool
}

func NewManager(pool *pgxpool.Pool) *Manager {
	return &Manager{pool: pool}
}

func (m *Manager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	if existingTx := ctx.Value(txKey{}); existingTx != nil {
		return fn(ctx)
	}

	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}

	ctxWithTx := context.WithValue(ctx, txKey{}, tx)

	err = fn(ctxWithTx)
	if err != nil {
		_ = tx.Rollback(ctxWithTx)
		return err
	}

	return tx.Commit(ctxWithTx)
}
