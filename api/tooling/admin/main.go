// This program performs administrative tasks for the bookstore service:
// schema migrations, demo data and user creation. It is built into the same
// image as the sales API so that the binary applying a migration and the
// binary serving traffic always come from the same commit.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/lmittmann/tint"
	"github.com/ousloob/bookstore/api/tooling/admin/command"
	"github.com/ousloob/bookstore/business/sys/database"
)

var build = "develop"

func main() {
	log := slog.New(tint.NewTextHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.DateTime,
	})).With("TOOL", "ADMIN")

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		if errors.Is(err, command.ErrHelp) {
			command.Usage(os.Stdout)
			return
		}

		log.ErrorContext(ctx, "admin", "msg", err)
		fmt.Fprintln(os.Stderr, `Run "admin help" for the list of commands.`)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *slog.Logger) error {

	// =========================================================================
	// Configuration
	//
	// The same DB_ variables the sales API reads, so the Kubernetes migration
	// job can mount the same secret. Only the settings this tool needs are
	// here: a one-shot command has no use for pool tuning.

	var cfg struct {
		DB struct {
			User       string `env:"USER"`
			Password   string `env:"PASSWORD"`
			Host       string `env:"HOST" envDefault:"localhost"`
			Name       string `env:"NAME" envDefault:"bookdb"`
			DisableTLS bool   `env:"DISABLE_TLS" envDefault:"true"`
		} `envPrefix:"DB_"`
	}

	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("parse conf: %w", err)
	}

	args := os.Args[1:]
	if len(args) == 0 {
		return command.ErrHelp
	}

	// =========================================================================
	// Commands that do not touch the database

	switch args[0] {
	case "help", "-h", "--help":
		return command.ErrHelp

	case "version":
		fmt.Println(build)
		return nil
	}

	// =========================================================================
	// Database support

	log.InfoContext(ctx, "startup", "status", "opening database pool", "host", cfg.DB.Host, "name", cfg.DB.Name)

	db, err := database.Open(ctx, database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	// =========================================================================
	// Commands

	switch args[0] {
	case "migrate":
		return command.Migrate(ctx, log, db, args[1:])

	case "seed":
		return command.Seed(ctx, log, db)

	case "useradd":
		return command.UserAdd(ctx, log, db, args[1:])

	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}
