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
	r.Handler(http.MethodGet, "/debug/health", web.Use(web.Handler{e, health}, middleware.RequestLogger(e.Log)))

	// signup/register routes
	r.Handler(http.MethodGet, "/signup", web.Use(web.Handler{e, signup}, middleware.RequestLogger(e.Log)))
	r.Handler(http.MethodGet, "/login", web.Use(web.Handler{e, login}, middleware.RequestLogger(e.Log)))

	// maybe routes
	mg := maybeGroup{
		maybe: maybe.New(e.Db),
	}
	r.Handler(http.MethodGet, "/", web.Use(web.Handler{e, mg.getAllMaybes}, middleware.RequestLogger(e.Log)))
	r.Handler(http.MethodGet, "/maybes", web.Use(web.Handler{e, mg.getMaybesQueryFilter}, middleware.RequestLogger(e.Log)))
	r.Handler(http.MethodGet, "/maybes/:id", web.Use(web.Handler{e, mg.getMaybeByID}, middleware.RequestLogger(e.Log)))

	// fileServer
	fileServer := http.FileServer(web.NeuteredFileSystem{http.Dir("./ui/static/")})
	r.Handler(http.MethodGet, "/static/*", http.StripPrefix("/static", fileServer))

	return r
}
