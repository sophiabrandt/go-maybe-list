package web

import (
	"encoding/json"
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
)

// Decode unmarhals an incoming JSON request.
func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return err
	}
	return nil
}

// Params returns the web call parameters from the request.
func Params(r *http.Request) map[string]string {
	return httptreemux.ContextParams(r.Context())
}
