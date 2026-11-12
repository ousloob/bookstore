package mid

import (
	"log/slog"
	"net/http"
	"time"
)

// Logger returns a middleware that logs one line per request, with the
// method, path, status, response size, duration and remote address.
func Logger(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := newResponseWriter(w)

			defer func() {
				log.InfoContext(r.Context(), "HTTP request",
					"method", r.Method,
					"path", r.URL.Path,
					"status", rw.Status(),
					"bytes", rw.BytesWritten(),
					"duration", time.Since(start),
					"remoteaddr", r.RemoteAddr,
				)
			}()

			next.ServeHTTP(rw, r)
		})
	}
}
