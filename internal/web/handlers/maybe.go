package handlers

import (
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/sophiabrandt/go-maybe-list/internal/data"
	"github.com/sophiabrandt/go-maybe-list/internal/data/maybe"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
	"github.com/sophiabrandt/go-maybe-list/internal/web/web"
)

type maybeGroup struct {
	maybe interface {
		Query() (maybe.Infos, error)
		QueryByID(id string) (maybe.Info, error)
		QueryByTitle(title string) (maybe.Infos, error)
	}
}

func (mg maybeGroup) getAllMaybes(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	maybes, err := mg.maybe.Query()
	if err != nil {
		return web.StatusError{Err: err, Code: http.StatusInternalServerError}
	}

	return web.Render(e, w, r, "home.page.tmpl", &data.TemplateData{Maybes: maybes}, http.StatusOK)
}

func (mg maybeGroup) getMaybesQueryFilter(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	title, err := url.QueryUnescape(r.URL.Query().Get("title"))
	if err != nil {
		return web.StatusError{Err: err, Code: http.StatusBadRequest}
	}
	maybes, err := mg.maybe.QueryByTitle(title)
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
	mb, err := mg.maybe.QueryByID(params["id"])
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

	return web.Render(e, w, r, "maybe_detail.page.tmpl", &data.TemplateData{Maybe: &mb}, http.StatusOK)
}
