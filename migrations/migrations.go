// Package migrations holds the database schema as SQL migrations embedded in
// the binary, and applies them. The seed data for demos lives in seed/ and is
// applied separately: migrations change the schema, seeds change the rows.
package migrations

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
)

//go:embed *.sql
var fsys embed.FS

//go:embed seed/*.sql
var seedFsys embed.FS

// ErrNoRollback is returned by Rollback when nothing has been applied yet.
var ErrNoRollback = errors.New("no migration to roll back")

// Status reports one migration and whether it has been applied.
type Status struct {
	Version   int64
	Source    string
	Applied   bool
	AppliedAt time.Time
}

// newProvider builds a goose provider over the embedded migrations. A lock
// table makes it safe to run from several instances at once: one migrates,
// the others wait.
//
// The returned close function releases the database/sql handle wrapping the
// pool. It does not close the pool itself.
func newProvider(log *slog.Logger, pool *pgxpool.Pool) (*goose.Provider, func() error, error) {
	db := stdlib.OpenDBFromPool(pool)

	locker, err := lock.NewPostgresTableLocker()
	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("creating migration locker: %w", err)
	}

	provider, err := goose.NewProvider(goose.DialectPostgres, db, fsys,
		goose.WithLocker(locker),
		goose.WithSlog(log),
	)
	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("creating migration provider: %w", err)
	}

	return provider, db.Close, nil
}

// Apply applies every pending migration.
func Apply(ctx context.Context, log *slog.Logger, pool *pgxpool.Pool) error {
	provider, closeDB, err := newProvider(log, pool)
	if err != nil {
		return err
	}
	defer closeDB()

	if _, err := provider.Up(ctx); err != nil {
		return fmt.Errorf("applying migrations: %w", err)
	}

	return nil
}

// Rollback rolls back the most recently applied migration, and only that one.
// Reversing a migration usually destroys data, so this exists for development.
// What keeps it out of production is the database role: the migrator role has
// the rights to run it, the application role does not.
func Rollback(ctx context.Context, log *slog.Logger, pool *pgxpool.Pool) (string, error) {
	provider, closeDB, err := newProvider(log, pool)
	if err != nil {
		return "", err
	}
	defer closeDB()

	result, err := provider.Down(ctx)
	if err != nil {
		if errors.Is(err, goose.ErrNoNextVersion) {
			return "", ErrNoRollback
		}

		return "", fmt.Errorf("rolling back migration: %w", err)
	}

	return result.Source.Path, nil
}

// Report lists every migration with the state the database is in.
func Report(ctx context.Context, log *slog.Logger, pool *pgxpool.Pool) ([]Status, error) {
	provider, closeDB, err := newProvider(log, pool)
	if err != nil {
		return nil, err
	}
	defer closeDB()

	results, err := provider.Status(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading migration status: %w", err)
	}

	statuses := make([]Status, 0, len(results))
	for _, result := range results {
		statuses = append(statuses, Status{
			Version:   result.Source.Version,
			Source:    result.Source.Path,
			Applied:   result.State == goose.StateApplied,
			AppliedAt: result.AppliedAt,
		})
	}

	return statuses, nil
}

// Seed loads the demo data. Every seed file is written to be idempotent, so
// running this twice is not an error, and all of them share one transaction:
// either the whole demo data set lands or none of it does.
func Seed(ctx context.Context, log *slog.Logger, pool *pgxpool.Pool) error {
	names, err := fs.Glob(seedFsys, "seed/*.sql")
	if err != nil {
		return fmt.Errorf("listing seed files: %w", err)
	}

	// fs.Glob sorts its results, so the files run in filename order.
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning seed transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, name := range names {
		stmts, err := fs.ReadFile(seedFsys, name)
		if err != nil {
			return fmt.Errorf("reading %s: %w", name, err)
		}

		log.InfoContext(ctx, "seeding", "file", name)

		if _, err := tx.Exec(ctx, string(stmts)); err != nil {
			return fmt.Errorf("running %s: %w", name, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing seed transaction: %w", err)
	}

	return nil
}
