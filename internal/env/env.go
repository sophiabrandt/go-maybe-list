package env

import (
	"log"

	"github.com/dimfeld/httptreemux/v5"
)

// Env defines the local app context and holds global
// dependencies.
type Env struct {
	Log    *log.Logger
	Router *httptreemux.ContextMux
}

// New creates a new pointer to an Env struct.
func New(log *log.Logger) *Env {
	router := httptreemux.NewContextMux()
	return &Env{
		Log:    log,
		Router: router,
	}
}
