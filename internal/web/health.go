package web

import (
	"net/http"

	"github.com/sophiabrandt/go-maybe-list/internal/env"
)

// health checks if the service is available and database is up.
func health(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	status := "ok"
	statusCode := http.StatusOK
	health := struct {
		Status string `json:"status`
	}{Status: status}

	return respond(e, w, health, statusCode)
}
