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
	user interface {
		QueryByID(userID string) (user.Info, error)
		Create(user user.NewUser) (user.Info, error)
		Authenticate(email, password string) (string, error)
		ChangePassword(currentPassword, newPassword, userID string) error
	}
}

func (ug userGroup) signupForm(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	return web.Render(e, w, r, "signup.page.tmpl", &data.TemplateData{Form: forms.New(nil)}, http.StatusOK)
}

func (ug userGroup) signup(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	// form validation
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRegex)
	form.MaxLength("name", 255)
	form.SecurePassword("password")
	form.IsEqualString("password", "confirm password")

	if !form.Valid() {
		return web.Render(e, w, r, "signup.page.tmpl", &data.TemplateData{Form: form}, http.StatusUnprocessableEntity)
	}

	nu := user.NewUser{
		Name:            form.Get("name"),
		Email:           form.Get("email"),
		Password:        form.Get("password"),
		PasswordConfirm: form.Get("password_confirm"),
	}
	_, err := ug.user.Create(nu)
	if err != nil {
		switch errors.Cause(err) {
		case user.ErrDuplicateEmail:
			form.Errors.Add("email", "Invalid email or email already in use")
			return web.Render(e, w, r, "signup.page.tmpl", &data.TemplateData{Form: form}, http.StatusUnprocessableEntity)
		default:
			return errors.Wrapf(err, "creating new user: %+v", nu)
		}
	}

	e.Session.Put(r, "flash", "Signup successful. Please log in.")
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
	return nil
}

func (ug userGroup) logout(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	e.Session.Remove(r, "authenticatedUserID")
	e.Session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func (ug userGroup) loginForm(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	return web.Render(e, w, r, "login.page.tmpl", &data.TemplateData{Form: forms.New(nil)}, http.StatusOK)
}

func (ug userGroup) login(e *env.Env, w http.ResponseWriter, r *http.Request) error {
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

	path := e.Session.PopString(r, "redirectPathAfterLogin")
	if path != "" {
		http.Redirect(w, r, path, http.StatusSeeOther)
		return nil
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func (ug userGroup) profile(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	params := web.Params(r)
	userID := e.Session.GetString(r, "authenticatedUserID")

	usr, err := ug.user.QueryByID(userID)
	if err != nil {
		switch errors.Cause(err) {
		case user.ErrNotFound:
			return web.StatusError{Err: err, Code: http.StatusNotFound}
		default:
			return errors.Wrapf(err, "ID : %s", params["id"])
		}
	}

	return web.Render(e, w, r, "profile.page.tmpl", &data.TemplateData{User: &usr}, http.StatusOK)
}

func (ug userGroup) changePasswordForm(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	return web.Render(e, w, r, "changepassword.page.tmpl", &data.TemplateData{Form: forms.New(nil)}, http.StatusOK)
}

func (ug userGroup) changePassword(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	form := forms.New(r.PostForm)
	form.Required("current password", "password", "confirm password")
	form.SecurePassword("password")
	form.IsEqualString("password", "confirm password")

	if !form.Valid() {
		return web.Render(e, w, r, "changepassword.page.tmpl", &data.TemplateData{Form: form}, http.StatusUnprocessableEntity)
	}

	userID := e.Session.GetString(r, "authenticatedUserID")

	err := ug.user.ChangePassword(form.Get("current password"), form.Get("password"), userID)
	if err != nil {
		switch errors.Cause(err) {
		case user.ErrNotFound:
			return web.StatusError{Err: err, Code: http.StatusNotFound}
		case user.ErrAuthenticationFailure:
			form.Errors.Add("current password", "Current password is incorrect")
			return web.Render(e, w, r, "changepassword.page.tmpl", &data.TemplateData{Form: form}, http.StatusUnprocessableEntity)
		default:
			return errors.Wrapf(err, "ID : %s", userID)
		}
	}

	e.Session.Put(r, "flash", "Password successfully updated!")

	http.Redirect(w, r, "/users/profile", http.StatusSeeOther)
	return nil
}
