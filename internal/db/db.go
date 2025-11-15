package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/linspacestrom/InterShipAv/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	tryConnectToDb           = 10
	firstIntervalToConnectDB = 500 * time.Millisecond
	intervalToPing           = 300 * time.Millisecond
	waitBeforePing           = 50 * time.Millisecond
)

func buildPostgresDSN(cfg config.DbConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
}

func NewPoolPostgres(ctx context.Context, cfg config.DbConfig) (*pgxpool.Pool, error) {
	log.Printf("Waiting %v before first database ping...\n", firstIntervalToConnectDB)
	select {
	case <-time.After(firstIntervalToConnectDB):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	dsn := buildPostgresDSN(cfg)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	for i := 1; i <= tryConnectToDb; i++ {
		pingCtx, cancel := context.WithTimeout(ctx, waitBeforePing)
		err := pool.Ping(pingCtx)
		cancel()

		if err == nil {
			log.Printf("database ping successful on attempt %d\n", i)
			return pool, nil
		}

		log.Printf("attempt %d: database not ready, retrying...\n", i)

		select {
		case <-time.After(intervalToPing):
		case <-ctx.Done():
			pool.Close()
			return nil, ctx.Err()
		}
	}

	pool.Close()
	return nil, fmt.Errorf("database not ready after %d attempts", tryConnectToDb)
}
