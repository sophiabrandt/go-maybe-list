package data

import "github.com/sophiabrandt/go-maybe-list/internal/data/maybe"

// TemplateData is all data needed for Golang HTML templates.
type TemplateData struct {
	Maybe       *maybe.Info
	Maybes      maybe.Infos
	Flash       string
	CurrentYear int
}
