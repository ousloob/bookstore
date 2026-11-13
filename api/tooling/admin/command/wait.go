package command

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// attemptTimeout bounds a single ping, so a database that accepts the
// connection and then stalls does not eat the whole budget in one attempt.
const attemptTimeout = 2 * time.Second

// Wait blocks until the database answers, or gives up. It runs as the init
// container of the migration job: migrate itself fails fast and loudly, so
// something else has to absorb the ordinary case of a database that is still
// starting up.
//
// It pings through the pool, which takes a connection from it, so success
// means the credentials work and the database exists. A port that merely
// accepts connections is not enough: PostgreSQL restarts itself partway
// through initdb, and it answers on the port before the role and the database
// are there.
func Wait(ctx context.Context, log *slog.Logger, db *pgxpool.Pool, args []string) error {
	fs := flag.NewFlagSet("wait", flag.ContinueOnError)
	timeout := fs.Duration("timeout", time.Minute, "how long to keep trying before giving up")
	interval := fs.Duration("interval", time.Second, "how long to wait between attempts")

	if err := fs.Parse(args); err != nil {
		return ErrHelp
	}

	ctx, cancel := context.WithTimeout(ctx, *timeout)
	defer cancel()

	for attempt := 1; ; attempt++ {
		err := ping(ctx, db)
		if err == nil {
			log.InfoContext(ctx, "wait", "status", "database ready", "attempts", attempt)
			return nil
		}

		log.InfoContext(ctx, "wait", "status", "database not ready yet", "attempt", attempt, "err", err)

		select {
		case <-ctx.Done():
			return fmt.Errorf("wait: database not ready after %s: %w", *timeout, err)

		case <-time.After(*interval):
		}
	}
}

func ping(ctx context.Context, db *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(ctx, attemptTimeout)
	defer cancel()

	return db.Ping(ctx)
}
