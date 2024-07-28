package web

import (
	"net/http"

	"github.com/sophiabrandt/go-maybe-list/internal/env"
)

// ParamByName returns the value of the first param with the given name.
func ParamByName(r *http.Request, name string) string {
	return r.PathValue(name)
}

// IsAuthenticated checks the current request for an authenticated user.
func IsAuthenticated(e *env.Env, r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(ContextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
