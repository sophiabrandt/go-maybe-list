package web

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
)

// Params returns the web call parameters from the request.
func Params(r *http.Request) httprouter.Params {
	return httprouter.ParamsFromContext(r.Context())
}

// IsAuthenticated checks the current request for an authenticated user.
func IsAuthenticated(e *env.Env, r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(ContextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
