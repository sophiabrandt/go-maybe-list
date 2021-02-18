// https://blog.questionable.services/article/http-handler-error-handling-revisited/
package web

import (
	"net/http"

	"github.com/sophiabrandt/go-maybe-list/internal/env"
)

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

// ErrorResponse is the client response struct for errors.
type ErrorResponse struct {
	Error  string   `json:"error"`
	Fields []string `json:"fields,omitempty"`
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

// handler takes a configured Env.
type handler struct {
	E *env.Env
	H func(E *env.Env, w http.ResponseWriter, r *http.Request) error
}

// ServeHTTP allows the Handler to satisy the http.Handler interface.
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.E, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			h.E.Log.Printf("HTTP %d - %s", e.Status(), e)
			response := ErrorResponse{e.Error(), nil}
			respond(h.E, w, response, e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			response := ErrorResponse{http.StatusText(http.StatusInternalServerError), nil}
			h.E.Log.Printf("%s", e)
			respond(h.E, w, response, http.StatusInternalServerError)
		}
	}
}

// use wraps middleware around handlers.
func use(h handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	var res http.Handler = h
	for _, m := range middleware {
		res = m(res)
	}

	return res
}
