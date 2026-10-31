package mid

import "net/http"

// responseRecorder wraps an http.ResponseWriter to capture the status code
// and the number of body bytes written, for logging.
type responseRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *responseRecorder) WriteHeader(status int) {
	if r.status == 0 {
		r.status = status
	}
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(b)
	r.bytes += n
	return n, err
}

// Unwrap lets http.ResponseController reach the underlying writer's
// optional interfaces (Flush, Hijack, SetWriteDeadline, ...).
func (r *responseRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

// Status returns the response status code, or 200 if the handler wrote a
// body without calling WriteHeader, or 0 if nothing was written at all.
func (r *responseRecorder) Status() int {
	return r.status
}
