package env

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// Env defines the local app context and holds global
// dependencies.
type Env struct {
	Log           *log.Logger
	TemplateCache map[string]*template.Template
	Session       *scs.SessionManager
}

// New creates a new pointer to an Env struct.
func New(log *log.Logger, templateCache map[string]*template.Template, session *scs.SessionManager) *Env {
	return &Env{
		Log:           log,
		TemplateCache: templateCache,
		Session:       session,
	}
}
