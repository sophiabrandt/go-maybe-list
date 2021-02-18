package env

import (
	"html/template"
	"log"

	"github.com/dimfeld/httptreemux/v5"
)

// Env defines the local app context and holds global
// dependencies.
type Env struct {
	Log           *log.Logger
	Router        *httptreemux.ContextMux
	TemplateCache map[string]*template.Template
}

// New creates a new pointer to an Env struct.
func New(log *log.Logger, templateCache map[string]*template.Template) *Env {
	router := httptreemux.NewContextMux()
	return &Env{
		Log:           log,
		Router:        router,
		TemplateCache: templateCache,
	}
}
