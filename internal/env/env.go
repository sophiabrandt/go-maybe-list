package env

import (
	"html/template"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/golangcollege/sessions"
)

// Env defines the local app context and holds global
// dependencies.
type Env struct {
	Log           *log.Logger
	Db            *sqlx.DB
	Router        *httptreemux.ContextMux
	TemplateCache map[string]*template.Template
	Validator     *validator.Validate
	Session       *sessions.Session
}

// New creates a new pointer to an Env struct.
func New(log *log.Logger, db *sqlx.DB, templateCache map[string]*template.Template, validator *validator.Validate, session *sessions.Session) *Env {
	router := httptreemux.NewContextMux()
	return &Env{
		Log:           log,
		Db:            db,
		Router:        router,
		TemplateCache: templateCache,
		Validator:     validator,
		Session:       session,
	}
}
