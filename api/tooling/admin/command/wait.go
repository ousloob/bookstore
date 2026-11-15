package command

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ousloob/bookstore/migrations"
)

// attemptTimeout bounds a single check, so a database that accepts the
// connection and then stalls does not eat the whole budget in one attempt.
const attemptTimeout = 2 * time.Second

// Wait blocks until the database is ready, or gives up. It is the init
// container of both pods that need the database:
//
//   - the migration job waits for the database to answer, so migrate itself
//     can fail fast and loudly at the first real error;
//   - the API runs `wait migrations`, which also waits until the schema has
//     every migration this binary carries, so it never serves against a
//     schema it does not know, whoever applied the Deployment and in
//     whatever order.
//
// It pings through the pool, which takes a connection from it, so success
// means the credentials work and the database exists. A port that merely
// accepts connections is not enough: PostgreSQL restarts itself partway
// through initdb, and it answers on the port before the role and the database
// are there.
func Wait(ctx context.Context, log *slog.Logger, db *pgxpool.Pool, args []string) error {
	if len(args) == 0 {
		return errors.New("wait needs a subcommand: database or migrations")
	}

	var migrated bool

	switch args[0] {
	case "database":
	case "migrations":
		migrated = true
	default:
		return fmt.Errorf("unknown wait subcommand %q, want database or migrations", args[0])
	}

	fs := flag.NewFlagSet("wait", flag.ContinueOnError)
	timeout := fs.Duration("timeout", time.Minute, "how long to keep trying before giving up")
	interval := fs.Duration("interval", time.Second, "how long to wait between attempts")

	if err := fs.Parse(args[1:]); err != nil {
		return ErrHelp
	}

	ctx, cancel := context.WithTimeout(ctx, *timeout)
	defer cancel()

	for attempt := 1; ; attempt++ {
		err := check(ctx, log, db, migrated)
		if err == nil {
			log.InfoContext(ctx, "wait", "status", "ready", "attempts", attempt)
			return nil
		}

		log.InfoContext(ctx, "wait", "status", "not ready yet", "attempt", attempt, "err", err)

		select {
		case <-ctx.Done():
			return fmt.Errorf("wait: not ready after %s: %w", *timeout, err)

		case <-time.After(*interval):
		}
	}
}

func check(ctx context.Context, log *slog.Logger, db *pgxpool.Pool, migrated bool) error {
	ctx, cancel := context.WithTimeout(ctx, attemptTimeout)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		return err
	}

	if !migrated {
		return nil
	}

	pending, err := migrations.Pending(ctx, log, db)
	if err != nil {
		return err
	}

	if pending {
		return errors.New("migrations pending, the migration job has not run yet")
	}

	return nil
}
