// Package v1 contains the full set of handler functions and routes supported
// by the v1 web api.
package v1

import (
	"log/slog"
	"net/http"
)

type Bus struct {
	Log *slog.Logger
}

// Routes binds all the version 1 routes.
func Routes(mux *http.ServeMux, bus Bus) {}
