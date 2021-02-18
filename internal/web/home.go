package web

import (
	"net/http"

	"github.com/sophiabrandt/go-maybe-list/internal/env"
	"github.com/sophiabrandt/go-maybe-list/internal/web/templates"
)

func index(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	if err := render(e, w, r, "home.page.tmpl", &templates.TemplateData{}); err != nil {
		return StatusError{err, http.StatusInternalServerError}
	}
	return nil
}

func signup(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	if err := render(e, w, r, "signup.page.tmpl", &templates.TemplateData{}); err != nil {
		return StatusError{err, http.StatusInternalServerError}
	}
	return nil
}

func login(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	if err := render(e, w, r, "login.page.tmpl", &templates.TemplateData{}); err != nil {
		return StatusError{err, http.StatusInternalServerError}
	}
	return nil
}
