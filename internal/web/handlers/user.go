package handlers

import (
	"net/http"

	"github.com/sophiabrandt/go-maybe-list/internal/data"
	"github.com/sophiabrandt/go-maybe-list/internal/data/user"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
	"github.com/sophiabrandt/go-maybe-list/internal/web/web"
)

type userGroup struct {
	user *user.RepositoryDb
}

func (ug userGroup) signupForm(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	return web.Render(e, w, r, "signup.page.tmpl", &data.TemplateData{}, http.StatusOK)
}

func (ug userGroup) loginForm(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	return web.Render(e, w, r, "signup.page.tmpl", &data.TemplateData{}, http.StatusOK)
}
