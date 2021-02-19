package handlers

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/sophiabrandt/go-maybe-list/internal/data"
	"github.com/sophiabrandt/go-maybe-list/internal/data/user"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
	"github.com/sophiabrandt/go-maybe-list/internal/web/forms"
	"github.com/sophiabrandt/go-maybe-list/internal/web/web"
)

type userGroup struct {
	user *user.RepositoryDb
}

func (ug userGroup) signupForm(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	return web.Render(e, w, r, "signup.page.tmpl", &data.TemplateData{Form: forms.New(nil)}, http.StatusOK)
}

func (ug userGroup) signup(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	var nu user.NewUser
	form, err := web.DecodeForm(r, &nu)
	if err != nil {
		return web.Render(e, w, r, "signup.page.tmpl", &data.TemplateData{Form: form}, http.StatusUnprocessableEntity)
	}

	_, err = ug.user.Create(nu)
	if err != nil {
		switch errors.Cause(err) {
		case user.ErrDuplicateEmail:
			form.Errors.Add("email", "Invalid email or email already in use")
			return web.Render(e, w, r, "signup.page.tmpl", &data.TemplateData{Form: form}, http.StatusUnprocessableEntity)
		default:
			return errors.Wrapf(err, "creating new user: %+v", nu)
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func (ug userGroup) logout(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	return web.Render(e, w, r, "signup.page.tmpl", &data.TemplateData{}, http.StatusOK)
}

func (ug userGroup) loginForm(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	return web.Render(e, w, r, "login.page.tmpl", &data.TemplateData{Form: forms.New(nil)}, http.StatusOK)
}

func (ug userGroup) login(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return web.StatusError{Err: err, Code: http.StatusBadRequest}
	}

	form := forms.New(r.PostForm)
	id, err := ug.user.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		switch errors.Cause(err) {
		case user.ErrAuthenticationFailure:
			form.Errors.Add("generic", "Email or Password is incorrect")
			return web.Render(e, w, r, "login.page.tmpl", &data.TemplateData{Form: form}, http.StatusUnprocessableEntity)
		default:
			return errors.Wrap(err, "authenticationg")
		}
	}

	e.Session.Put(r, "authenticatedUserID", id)

	http.Redirect(w, r, "/maybes", http.StatusSeeOther)
	return nil
}
