package handlers

import (
	"net/http"

	"github.com/sophiabrandt/go-maybe-list/internal/data/maybe"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
	"github.com/sophiabrandt/go-maybe-list/internal/web"
	"github.com/sophiabrandt/go-maybe-list/internal/web/middleware"
)

// New creates a new router with all application routes.
func New(e *env.Env) http.Handler {
	r := e.Router

	// liveness check
	r.Handler(http.MethodGet, "/debug/health", web.Use(web.Handler{e, Health}, middleware.RequestLogger(e.Log)))

	// signup/register routes
	r.Handler(http.MethodGet, "/signup", web.Use(web.Handler{e, Signup}, middleware.RequestLogger(e.Log)))
	r.Handler(http.MethodGet, "/login", web.Use(web.Handler{e, Login}, middleware.RequestLogger(e.Log)))

	// maybe routes
	mg := maybeGroup{
		maybe: maybe.New(e.Db),
	}
	r.Handler(http.MethodGet, "/", web.Use(web.Handler{e, mg.GetAllMaybes}, middleware.RequestLogger(e.Log)))

	// fileServer
	fileServer := http.FileServer(web.NeuteredFileSystem{http.Dir("./ui/static/")})
	r.Handler(http.MethodGet, "/static/*", http.StripPrefix("/static", fileServer))

	return r
}
