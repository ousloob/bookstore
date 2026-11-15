package command

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ousloob/bookstore/migrations"
)

// Migrate dispatches the migrate subcommands, named after the goose
// commands they run. The verb is required: "migrate" on its own used to
// mean up, and a word that changes the schema should be typed out.
func Migrate(ctx context.Context, log *slog.Logger, db *pgxpool.Pool, args []string) error {
	if len(args) == 0 {
		return errors.New("migrate needs a subcommand: up, down or status")
	}

	switch args[0] {
	case "up":
		return migrateUp(ctx, log, db)

	case "status":
		return migrateStatus(ctx, log, db)

	case "down":
		return migrateDown(ctx, log, db)

	default:
		return fmt.Errorf("unknown migrate subcommand %q, want up, down or status", args[0])
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
