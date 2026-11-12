// Package mid provides HTTP middleware and the Chain function that composes it.
package mid

import (
	"net/http"
	"slices"
)

// Middleware wraps a handler with extra behaviour before or after it runs.
type Middleware func(http.Handler) http.Handler

// Chain wraps h with the given middleware. The first middleware is the
// outermost one, so it runs first on the way in and last on the way out:
//
//	Chain(h, a, b) == a(b(h))
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for _, mw := range slices.Backward(mws) {
		h = mw(h)
	}
	return h
}
