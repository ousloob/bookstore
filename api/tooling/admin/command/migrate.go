package command

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ousloob/bookstore/migrations"
)

// Migrate dispatches the migrate subcommands. No subcommand means "up",
// because that is the one a deployment runs.
func Migrate(ctx context.Context, log *slog.Logger, db *pgxpool.Pool, args []string) error {
	sub := "up"
	if len(args) > 0 {
		sub = args[0]
	}

	switch sub {
	case "up":
		return migrateUp(ctx, log, db)

	case "status":
		return migrateStatus(ctx, log, db)

	case "down":
		return migrateDown(ctx, log, db)

	default:
		return fmt.Errorf("unknown migrate subcommand %q", sub)
	}
}

func migrateUp(ctx context.Context, log *slog.Logger, db *pgxpool.Pool) error {
	if err := migrations.Apply(ctx, log, db); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}

func migrateStatus(ctx context.Context, log *slog.Logger, db *pgxpool.Pool) error {
	statuses, err := migrations.Report(ctx, log, db)
	if err != nil {
		return fmt.Errorf("migrate status: %w", err)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "VERSION\tSTATE\tAPPLIED AT\tSOURCE")

	for _, status := range statuses {
		state, appliedAt := "pending", "-"
		if status.Applied {
			state = "applied"
			appliedAt = status.AppliedAt.Format(time.DateTime)
		}

		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\n",
			status.Version, state, appliedAt, filepath.Base(status.Source))
	}

	return tw.Flush()
}

func migrateDown(ctx context.Context, log *slog.Logger, db *pgxpool.Pool) error {
	if err := migrations.Rollback(ctx, log, db); err != nil {
		return fmt.Errorf("migrate down: %w", err)
	}

	return nil
}
