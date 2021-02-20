package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sophiabrandt/go-maybe-list/internal/adapter/database"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
)

type debugGroup struct {
	db *sqlx.DB
}

// health checks if the service is available and database is up.
func (dg debugGroup) health(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	status := "ok"
	statusCode := http.StatusOK
	if err := database.StatusCheck(dg.db); err != nil {
		status = "db not ready"
		statusCode = http.StatusInternalServerError
	}
	health := struct {
		Status string `json:"status`
	}{Status: status}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(health)
	return nil
}
