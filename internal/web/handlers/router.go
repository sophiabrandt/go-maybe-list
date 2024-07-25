package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/sophiabrandt/go-maybe-list/internal/data/maybe"
	"github.com/sophiabrandt/go-maybe-list/internal/data/user"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
	"github.com/sophiabrandt/go-maybe-list/internal/web/mid"
	"github.com/sophiabrandt/go-maybe-list/internal/web/web"
)

// New creates a new router with all application routes.
func New(e *env.Env, db *sqlx.DB) http.Handler {
	standardMiddleware := alice.New(mid.SecureHeaders, mid.LogRequest(e.Log), mid.RecoverPanic(e.Log))

	dynamicMiddleware := alice.New(e.Session.Enable, mid.NoSurf, mid.Authenticate(e, user.New(db)))

	r := httprouter.New()

	// liveness check
	dg := debugGroup{
		db: db,
	}
	r.Handler(http.MethodGet, "/debug/health", web.Handler{e, dg.health})

	// maybe routes
	mg := maybeGroup{
		maybe: maybe.New(db),
	}
	r.Handler(http.MethodGet, "/", dynamicMiddleware.Then(web.Handler{e, mg.getAllMaybes}))
	r.Handler(http.MethodGet, "/maybes/new", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{e, mg.createMaybeForm}))
	r.Handler(http.MethodPost, "/maybes/new", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{e, mg.createMaybe}))
	r.Handler(http.MethodGet, "/maybes/maybe/:id", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{e, mg.getMaybeByID}))
	r.Handler(http.MethodPost, "/maybes/maybe/:id/delete", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{e, mg.deleteMaybe}))
	r.Handler(http.MethodGet, "/maybes/maybe/:id/update", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{e, mg.updateMaybeForm}))
	r.Handler(http.MethodPost, "/maybes/maybe/:id/update", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{e, mg.updateMaybe}))
	r.Handler(http.MethodGet, "/tags", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{e, mg.getAllTags}))
	r.Handler(http.MethodGet, "/tags/tag/:id", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{e, mg.getMaybesByTag}))

	// user
	ug := userGroup{
		user: user.New(db),
	}
	r.Handler(http.MethodGet, "/users/signup", dynamicMiddleware.Then(web.Handler{e, ug.signupForm}))
	r.Handler(http.MethodPost, "/users/signup", dynamicMiddleware.Then(web.Handler{e, ug.signup}))
	r.Handler(http.MethodGet, "/users/login", dynamicMiddleware.Then(web.Handler{e, ug.loginForm}))
	r.Handler(http.MethodPost, "/users/login", dynamicMiddleware.Then(web.Handler{e, ug.login}))
	r.Handler(http.MethodPost, "/users/logout", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{e, ug.logout}))
	r.Handler(http.MethodGet, "/users/profile", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{e, ug.profile}))
	r.Handler(http.MethodGet, "/users/change-password", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{e, ug.changePasswordForm}))
	r.Handler(http.MethodPost, "/users/change-password", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{e, ug.changePassword}))

	// fileServer
	fileServer := http.FileServer(web.NeuteredFileSystem{http.Dir("./ui/static/")})
	r.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(r)
}
