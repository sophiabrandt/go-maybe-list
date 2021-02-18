package handlers

import (
	"net/http"

	"github.com/sophiabrandt/go-maybe-list/internal/data"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
	"github.com/sophiabrandt/go-maybe-list/internal/web"
)

func Index(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	if err := web.Render(e, w, r, "home.page.tmpl", &data.TemplateData{}); err != nil {
		return web.StatusError{err, http.StatusInternalServerError}
	}
	return nil
}

func Signup(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	if err := web.Render(e, w, r, "signup.page.tmpl", &data.TemplateData{}); err != nil {
		return web.StatusError{err, http.StatusInternalServerError}
	}
	return nil
}

func Login(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	if err := web.Render(e, w, r, "login.page.tmpl", &data.TemplateData{}); err != nil {
		return web.StatusError{err, http.StatusInternalServerError}
	}
	return nil
}
