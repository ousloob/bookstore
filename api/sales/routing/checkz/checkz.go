// Package checkz implements the Kubernetes probe endpoints. They are
// registered outside the application middleware chain: kubelet hits them
// every few seconds, which would drown the request log, so instead the
// readiness handler logs its own failures.
package checkz

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config carries the shared resources the probe handlers are built with.
type Config struct {
	Log *slog.Logger
	DB  *pgxpool.Pool
}

// Routes binds the probe endpoints, named after the ones the Kubernetes
// API server itself serves.
func Routes(mux *http.ServeMux, cfg Config) {
	mux.HandleFunc("GET /livez", liveness)
	mux.HandleFunc("GET /readyz", readiness(cfg.Log, cfg.DB))
}

// liveness reports that the process is up and serving.
func liveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "ok")
}

// readiness reports that the service can take traffic, which means the
// database answers.
func readiness(log *slog.Logger, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			log.ErrorContext(ctx, "readiness", "status", "database unreachable", "err", err)
			http.Error(w, "database unreachable", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ok")
	}
}
