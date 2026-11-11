package command

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ousloob/bookstore/migrations"
)

// Seed loads the demo data. It assumes the schema is current: run migrate
// first, which is what `make seed` does.
func Seed(ctx context.Context, log *slog.Logger, db *pgxpool.Pool) error {
	if err := migrations.Seed(ctx, log, db); err != nil {
		return fmt.Errorf("seed: %w", err)
	}

	log.InfoContext(ctx, "seed", "status", "demo data loaded")

	return nil
}
