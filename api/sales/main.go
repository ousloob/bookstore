package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/lmittmann/tint"
	"github.com/ousloob/bookstore/api/sales/routing"
	"github.com/ousloob/bookstore/business/sys/database"
)

var build = "develop"

func main() {
	logger := slog.New(tint.NewTextHandler(os.Stderr, &tint.Options{
		AddSource:  true,
		Level:      slog.LevelDebug,
		TimeFormat: time.DateTime,
	})).With("API", "SALES")

	ctx := context.Background()
	if err := run(ctx, logger); err != nil {
		logger.ErrorContext(ctx, "startup", "msg", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *slog.Logger) error {

	// =========================================================================
	// GOMAXPROCS

	log.InfoContext(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// =========================================================================
	// Configuration

	var cfg struct {
		Web struct {
			ReadTimeout     time.Duration `env:"READ_TIMEOUT" envDefault:"5s"`
			WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
			IdleTimeout     time.Duration `env:"IDLE_TIMEOUT" envDefault:"12s"`
			ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"20s"`
			APIHost         string        `env:"API_HOST" envDefault:"0.0.0.0:3000"`
		} `envPrefix:"WEB_"`
		DB struct {
			User                  string        `env:"USER"`
			Password              string        `env:"PASSWORD"`
			Host                  string        `env:"HOST" envDefault:"localhost"`
			Name                  string        `env:"NAME" envDefault:"bookdb"`
			MaxConns              int           `env:"MAX_CONNS" envDefault:"0"`
			MinConns              int           `env:"MIN_CONNS" envDefault:"0"`
			MinIdleConns          int           `env:"MIN_IDLE_CONNS" envDefault:"0"`
			MaxConnLifetime       time.Duration `env:"MAX_CONN_LIFETIME" envDefault:"0s"`
			MaxConnLifetimeJitter time.Duration `env:"MAX_CONN_LIFETIME_JITTER" envDefault:"0s"`
			MaxConnIdleTime       time.Duration `env:"MAX_CONN_IDLE_TIME" envDefault:"0s"`
			HealthCheckPeriod     time.Duration `env:"HEALTH_CHECK_PERIOD" envDefault:"0s"`
			PingTimeout           time.Duration `env:"PING_TIMEOUT" envDefault:"0s"`
			DisableTLS            bool          `env:"DISABLE_TLS" envDefault:"true"`
		} `envPrefix:"DB_"`
	}

	// =========================================================================
	// App Starting

	log.InfoContext(ctx, "starting service", "build", build)
	defer log.InfoContext(ctx, "shutdown complete")

	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("parse conf: %w", err)
	}

	// =========================================================================
	// Database Support

	log.InfoContext(ctx, "startup", "status", "initializing database support", "hostport", cfg.DB.Host)

	db, err := database.Open(ctx, database.Config{
		User:                  cfg.DB.User,
		Password:              cfg.DB.Password,
		Host:                  cfg.DB.Host,
		Name:                  cfg.DB.Name,
		MaxConns:              cfg.DB.MaxConns,
		MinConns:              cfg.DB.MinConns,
		MinIdleConns:          cfg.DB.MinIdleConns,
		MaxConnLifetime:       cfg.DB.MaxConnLifetime,
		MaxConnLifetimeJitter: cfg.DB.MaxConnLifetimeJitter,
		MaxConnIdleTime:       cfg.DB.MaxConnIdleTime,
		HealthCheckPeriod:     cfg.DB.HealthCheckPeriod,
		PingTimeout:           cfg.DB.PingTimeout,
		DisableTLS:            cfg.DB.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	// =========================================================================
	// Start API Service

	log.InfoContext(ctx, "startup", "status", "initializing V1 API support")

	// Make a channel to listen for an interrupt or terminal signal from the OS.
	// Use a buffered channel because the signal package require it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Construct the mux for the API calls.
	apiMux := routing.APIMux(
		routing.Config{
			Log:      log,
			DB:       db,
			Shutdown: shutdown,
		},
	)

	// Construct a server to service the request against the mux.
	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     slog.NewLogLogger(log.Handler(), slog.LevelError),
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this
	// error.

	serveErrors := make(chan error, 1)

	// Start the service listening for api requests.
	go func() {
		log.InfoContext(ctx, "startup", "status", "api router started", "host", api.Addr)
		serveErrors <- api.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown
	select {
	case err := <-serveErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.InfoContext(ctx, "shutdown", "status", "shutdown started", "signal", sig)

		// give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
