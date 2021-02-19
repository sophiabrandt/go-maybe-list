package handlers

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/sophiabrandt/go-maybe-list/internal/data/maybe"
	"github.com/sophiabrandt/go-maybe-list/internal/data/user"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
	"github.com/sophiabrandt/go-maybe-list/internal/web/middleware"
	"github.com/sophiabrandt/go-maybe-list/internal/web/web"
)

// New creates a new router with all application routes.
func New(e *env.Env) http.Handler {
	standardMiddleware := alice.New(middleware.SecureHeaders, middleware.LogRequest(e.Log), middleware.RecoverPanic(e.Log))

	dynamicMiddleware := alice.New(e.Session.Enable)

	r := e.Router

	// liveness check
	r.Handler(http.MethodGet, "/debug/health", web.Handler{e, health})

	// maybe routes
	mg := maybeGroup{
		maybe: maybe.New(e.Db),
	}
	r.Handler(http.MethodGet, "/", dynamicMiddleware.Then(web.Handler{e, mg.getAllMaybes}))
	r.Handler(http.MethodGet, "/maybes", dynamicMiddleware.Then(web.Handler{e, mg.getMaybesQueryFilter}))
	r.Handler(http.MethodGet, "/maybes/:id", dynamicMiddleware.Then(web.Handler{e, mg.getMaybeByID}))

	// user
	ug := userGroup{
		user: user.New(e.Db),
	}
	r.Handler(http.MethodGet, "/users/signup", dynamicMiddleware.Then(web.Handler{e, ug.signupForm}))
	r.Handler(http.MethodPost, "/users/signup", dynamicMiddleware.Then(web.Handler{e, ug.signup}))
	r.Handler(http.MethodGet, "/users/login", dynamicMiddleware.Then(web.Handler{e, ug.loginForm}))
	r.Handler(http.MethodPost, "/users/login", dynamicMiddleware.Then(web.Handler{e, ug.login}))
	r.Handler(http.MethodPost, "/users/logout", dynamicMiddleware.Then(web.Handler{e, ug.logout}))

	// fileServer
	fileServer := http.FileServer(web.NeuteredFileSystem{http.Dir("./ui/static/")})
	r.Handler(http.MethodGet, "/static/*", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(r)
}
