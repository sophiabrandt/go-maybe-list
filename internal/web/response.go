package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/sophiabrandt/go-maybe-list/internal/env"
)

// respond answers the client with JSON.
func respond(e *env.Env, w http.ResponseWriter, data interface{}, statusCode int) error {
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

// render renders a HTML page to the client.
func render(e *env.Env, w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) error {
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Frame-Options", "deny")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	ts, ok := e.TemplateCache[tmpl]
	if !ok {
		return StatusError{fmt.Errorf("Error: no HTML template for available for %s", tmpl), http.StatusInternalServerError}
	}

	buf := new(bytes.Buffer)
	err := ts.Execute(buf, r)
	if err != nil {
		return StatusError{err, http.StatusInternalServerError}
	}

	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
	return nil
}

// neuteredFileSystem is a custom file system to disable directory listings.
type neuteredFileSystem struct {
	fs http.FileSystem
}

// Open opens the files from the custom file system.
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
