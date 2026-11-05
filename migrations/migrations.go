// Package migrations holds the database schema as SQL migrations embedded in
// the binary, and applies them.
package migrations

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
)

//go:embed *.sql
var fsys embed.FS

// Apply applies every pending migration. A lock table makes it safe to run
// from several instances at once: one migrates, the others wait.
func Apply(ctx context.Context, log *slog.Logger, pool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	locker, err := lock.NewPostgresTableLocker()
	if err != nil {
		return fmt.Errorf("creating migration locker: %w", err)
	}

	provider, err := goose.NewProvider(goose.DialectPostgres, db, fsys,
		goose.WithLocker(locker),
		goose.WithSlog(log),
	)
	if err != nil {
		return fmt.Errorf("creating migration provider: %w", err)
	}

	if _, err := provider.Up(ctx); err != nil {
		return fmt.Errorf("applying migrations: %w", err)
	}

	return nil
}
