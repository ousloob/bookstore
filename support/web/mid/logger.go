package mid

import (
	"log/slog"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &responseRecorder{ResponseWriter: w}

		defer func() {
			slog.Info("HTTP request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.Status(),
				"bytes", rec.bytes,
				"duration", time.Since(start),
				"remoteaddr", r.RemoteAddr,
			)
		}()

		next.ServeHTTP(rec, r)
	})
}
