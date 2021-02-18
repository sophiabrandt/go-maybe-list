package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sophiabrandt/go-maybe-list/internal/data"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
)

// Respond answers the client with JSON.
func Respond(e *env.Env, w http.ResponseWriter, data interface{}, statusCode int) error {
	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Frame-Options", "deny")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}

// HumanDate returns time in YYYY-MM-DD format
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// addDefaultData adds data for all templates
func addDefaultData(dt *data.TemplateData, r *http.Request) *data.TemplateData {
	if dt == nil {
		dt = &data.TemplateData{}
	}
	dt.CurrentYear = time.Now().Year()

	return dt
}

// Render renders a HTML page to the client.
func Render(e *env.Env, w http.ResponseWriter, r *http.Request, tmpl string, dt interface{}, statusCode int) error {
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Frame-Options", "deny")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	switch d := dt.(type) {
	case *data.TemplateData:
		ts, ok := e.TemplateCache[tmpl]
		if !ok {
			return StatusError{fmt.Errorf("Error: no HTML template for available for %s", tmpl), http.StatusInternalServerError}
		}

		buf := new(bytes.Buffer)
		err := ts.Execute(buf, addDefaultData(d, r))
		if err != nil {
			return StatusError{err, http.StatusInternalServerError}
		}

		w.WriteHeader(statusCode)
		buf.WriteTo(w)

	case error:
		http.Error(w, d.Error(), statusCode)
		return nil

	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	return nil
}
