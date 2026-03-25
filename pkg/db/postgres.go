package db

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

var ErrNotConfigured = errors.New("application database is not configured")

type DB struct {
	pool *pgxpool.Pool
}

func Open(ctx context.Context, databaseURL string) (*DB, error) {
	databaseURL = strings.TrimSpace(databaseURL)
	if databaseURL == "" {
		return nil, nil
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create postgres pool")
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, errors.Wrap(err, "failed to ping postgres")
	}

	return &DB{pool: pool}, nil
}

func (d *DB) Close() {
	if d == nil || d.pool == nil {
		return
	}
	d.pool.Close()
}

func (d *DB) Pool() *pgxpool.Pool {
	if d == nil {
		return nil
	}
	return d.pool
}

func (d *DB) Migrate(ctx context.Context) error {
	if d == nil || d.pool == nil {
		return ErrNotConfigured
	}
	return ApplyMigrations(ctx, d.pool)
}
