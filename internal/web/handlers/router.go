package handlers

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/sophiabrandt/go-maybe-list/internal/data/maybe"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
	"github.com/sophiabrandt/go-maybe-list/internal/web"
	"github.com/sophiabrandt/go-maybe-list/internal/web/middleware"
)

// New creates a new router with all application routes.
func New(e *env.Env) http.Handler {
	standardMiddleware := alice.New(middleware.RecoverPanic, middleware.LogRequest(e.Log))

	dynamicMiddleware := alice.New(e.Session.Enable)

	r := e.Router

	// liveness check
	r.Handler(http.MethodGet, "/debug/health", web.Handler{e, health})

	// signup/register routes
	r.Handler(http.MethodGet, "/signup", dynamicMiddleware.Then(web.Handler{e, signup}))
	r.Handler(http.MethodGet, "/login", dynamicMiddleware.Then(web.Handler{e, login}))

	// maybe routes
	mg := maybeGroup{
		maybe: maybe.New(e.Db),
	}
	r.Handler(http.MethodGet, "/", dynamicMiddleware.Then(web.Handler{e, mg.getAllMaybes}))
	r.Handler(http.MethodGet, "/maybes", dynamicMiddleware.Then(web.Handler{e, mg.getMaybesQueryFilter}))
	r.Handler(http.MethodGet, "/maybes/:id", dynamicMiddleware.Then(web.Handler{e, mg.getMaybeByID}))

	// fileServer
	fileServer := http.FileServer(web.NeuteredFileSystem{http.Dir("./ui/static/")})
	r.Handler(http.MethodGet, "/static/*", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(r)
}
