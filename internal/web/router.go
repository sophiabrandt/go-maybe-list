package web

import (
	"net/http"

	"github.com/sophiabrandt/go-maybe-list/internal/env"
	"github.com/sophiabrandt/go-maybe-list/internal/web/middleware"
)

// NewRouter creates a new router with all application routes.
func NewRouter(e *env.Env) http.Handler {
	r := e.Router

	r.Handler(http.MethodGet, "/debug/health", use(handler{e, health}, middleware.RequestLogger(e.Log)))

	return r
}
