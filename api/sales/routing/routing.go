// Package routing manages the API routing and middleware setup.
package routing

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	v1 "github.com/ousloob/bookstore/api/sales/routing/v1"
	"github.com/ousloob/bookstore/support/web/mid"
)

// Config carries the shared resources the routes are built with.
type Config struct {
	Log      *slog.Logger
	DB       *pgxpool.Pool
	Shutdown chan os.Signal
}

// APIMux constructs and returns the API handler: the router wrapped in
// the application-wide middleware.
func APIMux(cfg Config) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthz)
	mux.HandleFunc("GET /readyz", readyz(cfg.DB))

	v1.Routes(mux, v1.Config{Log: cfg.Log})

	return mid.Chain(mux, mid.Logger(cfg.Log))
}

// healthz is the liveness probe: the process is up and serving.
func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "ok")
}

// readyz is the readiness probe: ready to take traffic, which means the
// database answers.
func readyz(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			http.Error(w, "database unreachable", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ok")
	}
}
