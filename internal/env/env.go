package env

import (
	"html/template"
	"log"

	"github.com/golangcollege/sessions"
)

// Env defines the local app context and holds global
// dependencies.
type Env struct {
	Log           *log.Logger
	TemplateCache map[string]*template.Template
	Session       *sessions.Session
}

// New creates a new pointer to an Env struct.
func New(log *log.Logger, templateCache map[string]*template.Template, session *sessions.Session) *Env {
	return &Env{
		Log:           log,
		TemplateCache: templateCache,
		Session:       session,
	}
}
