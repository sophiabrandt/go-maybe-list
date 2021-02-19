package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/form/v4"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
	"github.com/sophiabrandt/go-maybe-list/internal/web/forms"
)

// validate holds the settings and caches for validating request struct values.
var validate *validator.Validate

const (
	dateRegexString string = "^(((19|20)([2468][048]|[13579][26]|0[48])|2000)[/-]02[/-]29|((19|20)[0-9]{2}[/-](0[469]|11)[/-](0[1-9]|[12][0-9]|30)|(19|20)[0-9]{2}[/-](0[13578]|1[02])[/-](0[1-9]|[12][0-9]|3[01])|(19|20)[0-9]{2}[/-]02[/-](0[1-9]|1[0-9]|2[0-8])))$"
)

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

	validate.RegisterValidation("date", isDate)
	validate.RegisterValidation("secure_password", isSecurePassword)

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

// IsAuthenticated checks the current request for an authenticated user.
func IsAuthenticated(e *env.Env, r *http.Request) bool {
	return e.Session.Exists(r, "authenticatedUserID")
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
				form.Errors.Add(verror.Field(), fmt.Sprintf("%s must be equal to %s", verror.Field(), verror.Param()))
			case "secure_password":
				form.Errors.Add(verror.Field(), fmt.Sprintf("%s must be at least 8 characters, contain at least 1 uppercase letter, at least 1 lowercase letter, at least 1 special character and at least 1 number", verror.Field()))
			case "date":
				form.Errors.Add(verror.Field(), fmt.Sprintf("%s must be a valid date", verror.Field()))
			default:
				form.Errors.Add(verror.Field(), fmt.Sprintf("something wrong on %s; %s", verror.Field(), verror.Tag()))
			}
		}
		return form, err
	}

	return form, nil
}

func isDate(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile(dateRegexString)
	return reg.MatchString(fl.Field().String())
}

// http://www.inanzzz.com/index.php/post/8l1a/validating-user-password-in-golang-requests
func isSecurePassword(fl validator.FieldLevel) bool {
	var (
		upp, low, num, sym bool
		tot                uint8
	)

	for _, char := range fl.Field().String() {
		switch {
		case unicode.IsUpper(char):
			upp = true
			tot++
		case unicode.IsLower(char):
			low = true
			tot++
		case unicode.IsNumber(char):
			num = true
			tot++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
			tot++
		default:
			return false
		}
	}

	if !upp || !low || !num || !sym || tot < 8 {
		return false
	}

	return true
}
