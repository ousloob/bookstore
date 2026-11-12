// Package database provides support for accessing the database.
package database

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Set of error variables for CRUD operations.
var (
	ErrDuplicatedEntry = errors.New("duplicated entry")
)

// Config is the required properties to use the database. Every pool setting
// left at zero keeps the pgxpool default.
type Config struct {
	User     string
	Password string
	Host     string
	Name     string
	// MaxConns is the maximum number of connections in the pool. Default:
	// the greater of 4 and the number of CPUs.
	MaxConns int
	// MinConns is the minimum number of connections, idle or in use, the
	// pool keeps open. Default: 0.
	MinConns int
	// MinIdleConns is the number of idle connections the pool keeps open so
	// that a request does not have to wait for a new connection. Default: 0.
	MinIdleConns int
	// MaxConnLifetime is how long a connection may live before it is closed.
	// Default: 1h.
	MaxConnLifetime time.Duration
	// MaxConnLifetimeJitter adds a random duration up to this value to
	// MaxConnLifetime so connections are not all closed at once. Default: 0.
	MaxConnLifetimeJitter time.Duration
	// MaxConnIdleTime is how long a connection may stay idle before it is
	// closed by the health check. Default: 30m.
	MaxConnIdleTime time.Duration
	// HealthCheckPeriod is how often the pool checks idle connections.
	// Default: 1m.
	HealthCheckPeriod time.Duration
	// PingTimeout is how long the pool waits for a ping reply before it
	// treats the connection as broken. Default: no timeout.
	PingTimeout time.Duration
	DisableTLS  bool
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
	if cfg.MinConns > 0 {
		poolCfg.MinConns = int32(cfg.MinConns)
	}
	if cfg.MinIdleConns > 0 {
		poolCfg.MinIdleConns = int32(cfg.MinIdleConns)
	}
	if cfg.MaxConnLifetime > 0 {
		poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	}
	if cfg.MaxConnLifetimeJitter > 0 {
		poolCfg.MaxConnLifetimeJitter = cfg.MaxConnLifetimeJitter
	}
	if cfg.MaxConnIdleTime > 0 {
		poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	}
	if cfg.HealthCheckPeriod > 0 {
		poolCfg.HealthCheckPeriod = cfg.HealthCheckPeriod
	}
	if cfg.PingTimeout > 0 {
		poolCfg.PingTimeout = cfg.PingTimeout
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("creating connection pool: %w", err)
	}

	return pool, nil
}
