package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/sophiabrandt/go-maybe-list/internal/data"
	"github.com/sophiabrandt/go-maybe-list/internal/data/maybe"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
	"github.com/sophiabrandt/go-maybe-list/internal/web/forms"
	"github.com/sophiabrandt/go-maybe-list/internal/web/web"
)

type maybeGroup struct {
	maybe interface {
		Query(userID string) (maybe.Infos, error)
		QueryByID(maybeID, userID string) (maybe.Info, error)
		QueryByTitle(title, userID string) (maybe.Infos, error)
		QueryByTag(maybeID, userID string) (maybe.Infos, error)
		Create(nm maybe.NewOrUpdateMaybe, userID string) (maybe.Info, error)
		Update(um maybe.NewOrUpdateMaybe, maybeID string, userID string) error
		Delete(maybeID string) error
	}
}

func (mg maybeGroup) getAllMaybes(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	userID := e.Session.GetString(r, "authenticatedUserID")
	maybes, err := mg.maybe.Query(userID)
	if err != nil {
		switch errors.Cause(err) {
		case maybe.ErrNotFound:
			return web.StatusError{Err: err, Code: http.StatusNotFound}
		default:
			return web.StatusError{Err: err, Code: http.StatusInternalServerError}
		}
	}

	return web.Render(e, w, r, "home.page.tmpl", &data.TemplateData{Maybes: maybes}, http.StatusOK)
}

func (mg maybeGroup) getMaybesByTag(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	params := web.Params(r)
	userID := e.Session.GetString(r, "authenticatedUserID")

	maybes, err := mg.maybe.QueryByTag(params["id"], userID)
	if err != nil {
		switch errors.Cause(err) {
		case maybe.ErrInvalidTag:
			return web.StatusError{Err: err, Code: http.StatusBadRequest}
		case maybe.ErrNotFound:
			return web.StatusError{Err: err, Code: http.StatusNotFound}
		default:
			return web.StatusError{Err: err, Code: http.StatusInternalServerError}
		}
	}

	return web.Render(e, w, r, "home.page.tmpl", &data.TemplateData{Maybes: maybes}, http.StatusOK)
}

func (mg maybeGroup) getMaybesQueryFilter(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	title, err := url.QueryUnescape(r.URL.Query().Get("title"))
	userID := e.Session.GetString(r, "authenticatedUserID")
	if err != nil {
		return web.StatusError{Err: err, Code: http.StatusBadRequest}
	}
	maybes, err := mg.maybe.QueryByTitle(title, userID)
	if err != nil {
		switch errors.Cause(err) {
		case maybe.ErrNotFound:
			return web.StatusError{Err: err, Code: http.StatusNotFound}
		default:
			return errors.Wrapf(err, "Query Path: %s, Title: %s", r.URL.EscapedPath(), title)
		}
	}

	return web.Render(e, w, r, "maybe.page.tmpl", &data.TemplateData{Maybes: maybes}, http.StatusOK)
}

func (mg maybeGroup) getMaybeByID(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	params := web.Params(r)
	userID := e.Session.GetString(r, "authenticatedUserID")

	mb, err := mg.maybe.QueryByID(params["id"], userID)
	if err != nil {
		switch errors.Cause(err) {
		case maybe.ErrInvalidID:
			return web.StatusError{Err: err, Code: http.StatusBadRequest}
		case maybe.ErrForbidden:
			return web.StatusError{Err: err, Code: http.StatusForbidden}
		case maybe.ErrNotFound:
			return web.StatusError{Err: err, Code: http.StatusNotFound}
		default:
			return errors.Wrapf(err, "ID : %s", params["id"])
		}
	}

	return web.Render(e, w, r, "maybe_detail.page.tmpl", &data.TemplateData{Maybe: &mb}, http.StatusOK)
}

func (mg maybeGroup) createMaybeForm(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	return web.Render(e, w, r, "create.page.tmpl", &data.TemplateData{Form: forms.New(nil)}, http.StatusOK)
}

func (mg maybeGroup) createMaybe(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	// form validation
	form := forms.New(r.PostForm)
	form.Required("title", "url", "description")
	form.ValidUrl("url")
	form.MaxLength("title", 255)
	form.MaxLength("description", 255)

	if !form.Valid() {
		return web.Render(e, w, r, "create.page.tmpl", &data.TemplateData{Form: form}, http.StatusUnprocessableEntity)
	}

	nm := maybe.NewOrUpdateMaybe{
		Title:       strings.TrimSpace(form.Get("title")),
		Url:         form.Get("url"),
		Description: form.Get("description"),
		Tags:        nil,
	}

	// if user added tags into the form, make a slice of tags,
	// sanitize them and add them to the model
	tags := strings.TrimSpace(form.Get("tags"))
	if tags != "" {
		t := strings.Split(tags, ",")
		// trim whitespace
		trimT := make([]string, len(t))
		for i, tag := range t {
			trimT[i] = strings.TrimSpace(tag)
		}
		nm.Tags = trimT
	}

	userID := e.Session.GetString(r, "authenticatedUserID")

	myb, err := mg.maybe.Create(nm, userID)
	if err != nil {
		switch errors.Cause(err) {
		case maybe.ErrInvalidTag:
			form.Errors.Add("tags", "tags are invalid")
			return web.Render(e, w, r, "create.page.tmpl", &data.TemplateData{Form: form}, http.StatusUnprocessableEntity)
		default:
			return errors.Wrapf(err, "creating new maybe: %v", nm)
		}
	}

	e.Session.Put(r, "flash", "Maybe successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/maybes/%v", myb.ID), http.StatusSeeOther)

	return nil
}

func (mg maybeGroup) updateMaybeForm(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	params := web.Params(r)
	userID := e.Session.GetString(r, "authenticatedUserID")

	mb, err := mg.maybe.QueryByID(params["id"], userID)
	if err != nil {
		switch errors.Cause(err) {
		case maybe.ErrInvalidID:
			return web.StatusError{Err: err, Code: http.StatusBadRequest}
		case maybe.ErrNotFound:
			return web.StatusError{Err: err, Code: http.StatusNotFound}
		default:
			return errors.Wrapf(err, "ID : %s", params["id"])
		}
	}

	// populate form with previous values
	form := forms.New(url.Values{})
	form.Set("title", mb.Title)
	form.Set("url", mb.Url)
	form.Set("description", mb.Description)
	// convert tag model into a list of strings for the form
	if mb.Tags != nil {
		tagNames := make([]string, len(mb.Tags))
		for i, tag := range mb.Tags {
			tagNames[i] = tag.Name
		}
		form.Set("tags", strings.Join(tagNames[:], ","))
	}

	return web.Render(e, w, r, "update.page.tmpl", &data.TemplateData{Maybe: &mb, Form: form}, http.StatusOK)
}

func (mg maybeGroup) updateMaybe(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	params := web.Params(r)
	// form validation
	form := forms.New(r.PostForm)
	form.MaxLength("title", 255)
	form.ValidUrl("url")
	form.MaxLength("description", 255)

	if !form.Valid() {
		return web.Render(e, w, r, "update.page.tmpl", &data.TemplateData{Form: form}, http.StatusUnprocessableEntity)
	}

	um := maybe.NewOrUpdateMaybe{
		Title:       strings.TrimSpace(form.Get("title")),
		Url:         form.Get("url"),
		Description: form.Get("description"),
		Tags:        nil,
	}

	// if user added tags into the form, make a slice of tags,
	// sanitize them and add them to the model
	tags := strings.TrimSpace(form.Get("tags"))
	if tags != "" {
		t := strings.Split(tags, ",")
		// trim whitespace
		trimT := make([]string, len(t))
		for i, tag := range t {
			trimT[i] = strings.TrimSpace(tag)
		}
		um.Tags = trimT
	}

	userID := e.Session.GetString(r, "authenticatedUserID")

	err := mg.maybe.Update(um, params["id"], userID)
	if err != nil {
		switch errors.Cause(err) {
		case maybe.ErrInvalidID:
			return web.StatusError{Err: err, Code: http.StatusBadRequest}
		case maybe.ErrForbidden:
			return web.StatusError{Err: err, Code: http.StatusForbidden}
		case maybe.ErrNotFound:
			return web.StatusError{Err: err, Code: http.StatusNotFound}
		default:
			return errors.Wrapf(err, "updating new maybe: %v", um)
		}
	}

	e.Session.Put(r, "flash", "Maybe successfully updated!")

	http.Redirect(w, r, fmt.Sprintf("/maybes/%v", params["id"]), http.StatusSeeOther)

	return nil
}

func (mg maybeGroup) deleteMaybe(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	params := web.Params(r)

	err := mg.maybe.Delete(params["id"])
	if err != nil {
		switch errors.Cause(err) {
		case maybe.ErrInvalidID:
			return web.StatusError{Err: err, Code: http.StatusBadRequest}
		case maybe.ErrNotFound:
			return web.StatusError{Err: err, Code: http.StatusNotFound}
		default:
			return errors.Wrapf(err, "deleting maybe with ID: %s", params["id"])
		}
	}

	e.Session.Put(r, "flash", "Maybe successfully deleted!")

	http.Redirect(w, r, "/maybes", http.StatusSeeOther)
	return nil
}
