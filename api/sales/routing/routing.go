// Package routing manages the API routing and middleware setup.
package routing

import (
	"log/slog"
	"net/http"
	"os"

	v1 "github.com/ousloob/bookstore/api/sales/routing/v1"
	"github.com/ousloob/bookstore/support/web/mid"
)

// Bus carries shared resources across the application.
type Bus struct {
	Log      *slog.Logger
	Shutdown chan os.Signal
}

// APIMux constructs and returns the API handler.
func APIMux(bus Bus) http.Handler {
	mux := http.NewServeMux()

	v1.Routes(mux, v1.Bus{Log: bus.Log})

	return mid.Logger(mux)
}
