package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/form/v4"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sophiabrandt/go-maybe-list/internal/web/forms"
)

// validate holds the settings and caches for validating request struct values.
var validate *validator.Validate

// NewValidator returns a pointer to a validator.
func NewValidator() *validator.Validate {
	// Instantiate the validator for use.
	validate = validator.New()

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validate
}

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

// DecodeForm parses incoming form data and validates it.
func DecodeForm(r *http.Request, val interface{}) (*forms.Form, error) {
	decoder := form.NewDecoder()
	// Use JSON tag names for errors instead of Go struct names.
	decoder.SetTagName("json")

	err := r.ParseForm()
	form := forms.New(r.PostForm)

	if err != nil {
		return form, StatusError{Err: err, Code: http.StatusBadRequest}
	}
	if err := decoder.Decode(&val, r.Form); err != nil {
		return form, StatusError{Err: err, Code: http.StatusBadRequest}
	}
	if err := validate.Struct(val); err != nil {

		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return form, err
		}

		// add the validation errors as form errors
		for _, verror := range verrors {
			switch verror.Tag() {
			case "required":
				form.Errors.Add(verror.Field(), fmt.Sprintf("%s is a required field", verror.Field()))
			case "min":
				form.Errors.Add(verror.Field(), fmt.Sprintf("%s must be a minimum of %s in length", verror.Field(), verror.Param()))
			case "max":
				form.Errors.Add(verror.Field(), fmt.Sprintf("%s must be a maximum of %s in length", verror.Field(), verror.Param()))
			case "url":
				form.Errors.Add(verror.Field(), fmt.Sprintf("%s must be a valid URL", verror.Field()))
			case "email":
				form.Errors.Add(verror.Field(), fmt.Sprintf("%s must be a valid Email", verror.Field()))
			case "eqfield":
				form.Errors.Add(verror.Field(), fmt.Sprintf("%s must be equal to %s", verror.Field(), verror.Tag()))
			default:
				form.Errors.Add(verror.Field(), fmt.Sprintf("something wrong on %s; %s", verror.Field(), verror.Tag()))
			}
		}
		return form, err
	}

	return form, nil
}
