package web

import (
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
)

// Params returns the web call parameters from the request.
func Params(r *http.Request) map[string]string {
	return httptreemux.ContextParams(r.Context())
}

// IsAuthenticated checks the current request for an authenticated user.
func IsAuthenticated(e *env.Env, r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(ContextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
