package db

import (
	"context"
	"embed"
	"io/fs"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type MigrationStatus struct {
	Name    string
	Applied bool
}

func MigrationNames() ([]string, error) {
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return nil, errors.Wrap(err, "failed to list embedded migrations")
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".sql") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

func ensureSchemaMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return ErrNotConfigured
	}
	_, err := pool.Exec(ctx, `
create table if not exists schema_migrations (
	name text primary key,
	applied_at timestamptz not null default now()
)`)
	if err != nil {
		return errors.Wrap(err, "failed to ensure schema_migrations table")
	}
	return nil
}

func ApplyMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return ErrNotConfigured
	}

	if err := ensureSchemaMigrations(ctx, pool); err != nil {
		return err
	}

	names, err := MigrationNames()
	if err != nil {
		return err
	}

	for _, name := range names {
		var applied bool
		if err := pool.QueryRow(ctx, "select exists(select 1 from schema_migrations where name = $1)", name).Scan(&applied); err != nil {
			return errors.Wrapf(err, "failed to query schema_migrations for %s", name)
		}
		if applied {
			continue
		}

		contents, err := migrationsFS.ReadFile("migrations/" + name)
		if err != nil {
			return errors.Wrapf(err, "failed to read migration %s", name)
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to begin migration transaction for %s", name)
		}

		if _, err := tx.Exec(ctx, string(contents)); err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrapf(err, "failed to execute migration %s", name)
		}

		if _, err := tx.Exec(ctx, "insert into schema_migrations(name) values($1)", name); err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrapf(err, "failed to record migration %s", name)
		}

		if err := tx.Commit(ctx); err != nil {
			return errors.Wrapf(err, "failed to commit migration %s", name)
		}
	}

	return nil
}

func MigrationStatuses(ctx context.Context, pool *pgxpool.Pool) ([]MigrationStatus, error) {
	if pool == nil {
		return nil, ErrNotConfigured
	}

	if err := ensureSchemaMigrations(ctx, pool); err != nil {
		return nil, err
	}

	names, err := MigrationNames()
	if err != nil {
		return nil, err
	}

	statuses := make([]MigrationStatus, 0, len(names))
	for _, name := range names {
		var applied bool
		if err := pool.QueryRow(ctx, "select exists(select 1 from schema_migrations where name = $1)", name).Scan(&applied); err != nil {
			return nil, errors.Wrapf(err, "failed to query migration status for %s", name)
		}
		statuses = append(statuses, MigrationStatus{Name: name, Applied: applied})
	}

	return statuses, nil
}
