package server

import (
	"net/http"
	"time"
)

// New creates a new http server.
func New(address string, handler http.Handler) *http.Server {
	srv := http.Server{
		Addr:         address,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}
	return &srv
}
