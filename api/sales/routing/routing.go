// Package routing manages the API routing and middleware setup.
package routing

import (
	"log/slog"
	"net/http"
	"os"

	v1 "github.com/ousloob/bookstore/api/sales/routing/v1"
	"github.com/ousloob/bookstore/support/web/mid"
)

// Config carries the shared resources the routes are built with.
type Config struct {
	Log      *slog.Logger
	Shutdown chan os.Signal
}

// APIMux constructs and returns the API handler: the router wrapped in
// the application-wide middleware.
func APIMux(cfg Config) http.Handler {
	mux := http.NewServeMux()

	v1.Routes(mux, v1.Config{Log: cfg.Log})

	return mid.Chain(mux, mid.Logger(cfg.Log))
}
