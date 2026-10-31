// Package database provides support for accessing the database.
package database

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Set of error variables for CRUD operations.
var (
	ErrDBDuplicatedEntry = errors.New("duplicated entry")
)

// Config is the required properties to use the database.
type Config struct {
	User     string
	Password string
	Host     string
	Name     string
	// MaxConns is the maximum number of connections in the pool. Zero keeps
	// the pgxpool default (the greater of 4 and the number of CPUs).
	MaxConns int
	// MinIdleConns is the number of idle connections the pool keeps open so
	// that a request does not have to wait for a new connection. Zero keeps
	// none.
	MinIdleConns int
	DisableTLS   bool
}

// Open knows how to open a database connection pool based on the
// configuration. It does not connect: the first query does.
func Open(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	poolCfg, err := pgxpool.ParseConfig(u.String())
	if err != nil {
		return nil, fmt.Errorf("parsing database config: %w", err)
	}

	if cfg.MaxConns > 0 {
		poolCfg.MaxConns = int32(cfg.MaxConns)
	}
	if cfg.MinIdleConns > 0 {
		poolCfg.MinIdleConns = int32(cfg.MinIdleConns)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("creating connection pool: %w", err)
	}

	return pool, nil
}
