package handlers

import (
	"net/http"

	"github.com/sophiabrandt/go-maybe-list/internal/adapter/database"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
	"github.com/sophiabrandt/go-maybe-list/internal/web"
)

// Health checks if the service is available and database is up.
func Health(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	status := "ok"
	statusCode := http.StatusOK
	if err := database.StatusCheck(e.Db); err != nil {
		status = "db not ready"
		statusCode = http.StatusInternalServerError
	}
	health := struct {
		Status string `json:"status`
	}{Status: status}

	return web.Respond(e, w, health, statusCode)
}
