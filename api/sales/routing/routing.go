// Package routing manages the API routing and middleware setup.
package routing

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ousloob/bookstore/api/sales/routing/checkz"
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
// the application-wide middleware. The probe endpoints sit outside the
// middleware chain; see package checkz for why.
func APIMux(cfg Config) http.Handler {
	api := http.NewServeMux()
	v1.Routes(api, v1.Config{Log: cfg.Log})

	mux := http.NewServeMux()
	checkz.Routes(mux, checkz.Config{Log: cfg.Log, DB: cfg.DB})
	mux.Handle("/", mid.Chain(api, mid.Logger(cfg.Log)))

	return mux
}
