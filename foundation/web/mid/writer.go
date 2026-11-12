package mid

import "net/http"

// responseWriter wraps an http.ResponseWriter and records the status code
// and the number of bytes written, so that middleware can log them after
// the handler has run.
type responseWriter struct {
	http.ResponseWriter
	status      int
	bytes       int
	wroteHeader bool
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	// A handler that writes a body without calling WriteHeader gets 200.
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

// WriteHeader records the status code. Only the first call counts, which
// mirrors net/http, where later calls are ignored.
func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

// Write records the number of bytes written.
func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.wroteHeader = true
	n, err := rw.ResponseWriter.Write(b)
	rw.bytes += n
	return n, err
}

// Unwrap exposes the underlying writer so that http.ResponseController can
// still reach Flush, Hijack and the deadline setters.
func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

// Status returns the status code sent to the client.
func (rw *responseWriter) Status() int {
	return rw.status
}

// BytesWritten returns the number of body bytes sent to the client.
func (rw *responseWriter) BytesWritten() int {
	return rw.bytes
}
