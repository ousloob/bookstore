// Package v1 contains the full set of handler functions and routes supported
// by the v1 web api.
package v1

import (
	"log/slog"
	"net/http"
)

// Config carries the shared resources the v1 routes are built with.
type Config struct {
	Log *slog.Logger
}

// Routes binds all the version 1 routes.
func Routes(mux *http.ServeMux, cfg Config) {}
