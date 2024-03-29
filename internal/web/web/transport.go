// https://blog.questionable.services/article/http-handler-error-handling-revisited/
package web

import (
	"net/http"
	"path/filepath"

	"github.com/sophiabrandt/go-maybe-list/internal/env"
)

type contextKey string

const ContextKeyIsAuthenticated = contextKey("isAuthenticated")

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Err  error
	Code int
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

// Handler takes a configured Env.
type Handler struct {
	E *env.Env
	H func(E *env.Env, w http.ResponseWriter, r *http.Request) error
}

// ServeHTTP allows the Handler to satisfy the http.Handler interface.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.E, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			h.E.Log.Printf("HTTP %d - %s", e.Status(), e)
			Render(h.E, w, r, "", e, e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			h.E.Log.Printf("HTTP 500 - %s", e)
			Render(h.E, w, r, "", http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

// NeuteredFileSystem is a custom file system to disable directory listings.
type NeuteredFileSystem struct {
	Fs http.FileSystem
}

// Open opens the files from the custom file system.
func (nfs NeuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.Fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.Fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
