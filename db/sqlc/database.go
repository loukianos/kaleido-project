package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Open connects a pool with pgx's default sizing and lifetimes. Pool tuning,
// when needed, belongs in DATABASE_URL (e.g. ?pool_max_conns=10).
func Open(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	return pool, nil
}

type Checker struct {
	DB *pgxpool.Pool
}

func (c Checker) Check(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return c.DB.Ping(ctx)
}
