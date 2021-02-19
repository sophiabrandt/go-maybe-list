package handlers

import (
	"net/http"

	"github.com/sophiabrandt/go-maybe-list/internal/data"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
	"github.com/sophiabrandt/go-maybe-list/internal/web/web"
)

func signup(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	return web.Render(e, w, r, "signup.page.tmpl", &data.TemplateData{}, http.StatusOK)
}

func login(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	return web.Render(e, w, r, "login.page.tmpl", &data.TemplateData{}, http.StatusOK)
}
