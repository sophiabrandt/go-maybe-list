package env

import (
	"html/template"
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/dimfeld/httptreemux/v5"
)

// Env defines the local app context and holds global
// dependencies.
type Env struct {
	Log           *log.Logger
	Db            *sqlx.DB
	Router        *httptreemux.ContextMux
	TemplateCache map[string]*template.Template
}

// New creates a new pointer to an Env struct.
func New(log *log.Logger, db *sqlx.DB, templateCache map[string]*template.Template) *Env {
	router := httptreemux.NewContextMux()
	return &Env{
		Log:           log,
		Db:            db,
		Router:        router,
		TemplateCache: templateCache,
	}
}
