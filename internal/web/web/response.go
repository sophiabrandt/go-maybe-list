package web

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/sophiabrandt/go-maybe-list/internal/data"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
)

// HumanDate returns time in YYYY-MM-DD format
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// addDefaultData adds data for all templates
func addDefaultData(e *env.Env, r *http.Request, dt *data.TemplateData) *data.TemplateData {
	if dt == nil {
		dt = &data.TemplateData{}
	}
	dt.CurrentYear = time.Now().Year()
	dt.Flash = e.Session.PopString(r, "flash")
	dt.IsAuthenticated = IsAuthenticated(e, r)

	return dt
}

// Render renders a HTML page to the client.
func Render(e *env.Env, w http.ResponseWriter, r *http.Request, tmpl string, dt interface{}, statusCode int) error {
	switch d := dt.(type) {
	case *data.TemplateData:
		ts, ok := e.TemplateCache[tmpl]
		if !ok {
			return StatusError{fmt.Errorf("Error: no HTML template for available for %s", tmpl), http.StatusInternalServerError}
		}

		buf := new(bytes.Buffer)
		err := ts.Execute(buf, addDefaultData(e, r, d))
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
