package handlers

import (
	"net/http"

	"github.com/sophiabrandt/go-maybe-list/internal/data"
	"github.com/sophiabrandt/go-maybe-list/internal/data/maybe"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
	"github.com/sophiabrandt/go-maybe-list/internal/web"
)

type maybeGroup struct {
	maybe *maybe.RepositoryDb
}

func (mg maybeGroup) GetAllMaybes(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	maybes, err := mg.maybe.Query()
	if err != nil {
		web.Render(e, w, r, "", err)
		return err
	}

	web.Render(e, w, r, "home.page.tmpl", &data.TemplateData{Maybes: maybes})
	return nil
}
