package data

import (
	"github.com/sophiabrandt/go-maybe-list/internal/data/maybe"
	"github.com/sophiabrandt/go-maybe-list/internal/web/forms"
)

// TemplateData is all data needed for Golang HTML templates.
type TemplateData struct {
	Maybe           *maybe.Info
	Maybes          maybe.Infos
	Form            *forms.Form
	Flash           string
	CurrentYear     int
	IsAuthenticated bool
	CSRFToken       string
}
