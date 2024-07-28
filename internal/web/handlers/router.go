package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"
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

	r := http.NewServeMux()

	// liveness check
	dg := debugGroup{
		db: db,
	}
	r.Handle("GET /debug/health", web.Handler{E: e, H: dg.health})

	// maybe routes
	mg := maybeGroup{
		maybe: maybe.New(db),
	}
	r.Handle("GET /", dynamicMiddleware.Then(web.Handler{E: e, H: mg.getAllMaybes}))
	r.Handle("GET /maybes/new", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{E: e, H: mg.createMaybeForm}))
	r.Handle("POST /maybes/new", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{E: e, H: mg.createMaybe}))
	r.Handle("GET /maybes/maybe/{id}", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{E: e, H: mg.getMaybeByID}))
	r.Handle("POST /maybes/maybe/{id}/delete", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{E: e, H: mg.deleteMaybe}))
	r.Handle("GET /maybes/maybe/{id}/update", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{E: e, H: mg.updateMaybeForm}))
	r.Handle("POST /maybes/maybe/{id}/update", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{E: e, H: mg.updateMaybe}))
	r.Handle("GET /tags", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{E: e, H: mg.getAllTags}))
	r.Handle("GET /tags/tag/{id}", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{E: e, H: mg.getMaybesByTag}))

	// user
	ug := userGroup{
		user: user.New(db),
	}
	r.Handle("GET /users/signup", dynamicMiddleware.Then(web.Handler{E: e, H: ug.signupForm}))
	r.Handle("POST /users/signup", dynamicMiddleware.Then(web.Handler{E: e, H: ug.signup}))
	r.Handle("GET /users/login", dynamicMiddleware.Then(web.Handler{E: e, H: ug.loginForm}))
	r.Handle("POST /users/login", dynamicMiddleware.Then(web.Handler{E: e, H: ug.login}))
	r.Handle("POST /users/logout", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{E: e, H: ug.logout}))
	r.Handle("GET /users/profile", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{E: e, H: ug.profile}))
	r.Handle("GET /users/change-password", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{E: e, H: ug.changePasswordForm}))
	r.Handle("POST /users/change-password", dynamicMiddleware.Append(mid.RequireAuthentication(e)).Then(web.Handler{E: e, H: ug.changePassword}))

	// fileServer
	fileServer := http.FileServer(web.NeuteredFileSystem{Fs: http.Dir("./ui/static/")})
	r.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(r)
}
