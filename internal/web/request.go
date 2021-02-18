package web

import (
	"encoding/json"
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
)

// decode unmarhals an incoming JSON request.
func decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return err
	}
	return nil
}

// params returns the web call parameters from the request.
func params(r *http.Request) map[string]string {
	return httptreemux.ContextParams(r.Context())
}
