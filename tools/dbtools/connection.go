package dbtools

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	errors "github.com/task_platform/tools/helpers"
)

// NewClient initializes the database connection pool.
func NewClient(ctx context.Context, dbURI string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dbURI)
	if err != nil {
		return nil, errors.Wrap("failed to parse DSN to config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap("unable to create connection pool: %w", err)
	}

	return pool, nil
}
