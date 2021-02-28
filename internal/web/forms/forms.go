// Package forms defines how to create a form and adds validations.
package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// MinLength checks for minimum length of characters.
func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", d))
	}
}

// MaxLength checks for maximum allowed length of characters.
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}

// PermittedValues disallows unknown fields.
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

// Required checks for empty fields.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field is required")
		}
	}
}

// Implement a MatchesPattern method to check that a specific field in the form
// matches a regular expression. If the check fails then add the
// appropriate message to the form errors.
func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}

// ValidUrl parses a field for a valid URL.
func (f *Form) ValidUrl(field string) {
	value := f.Get(field)
	_, err := url.ParseRequestURI(value)
	if err != nil {
		f.Errors.Add(field, "Invalid URL")
	}
}

// IsEqual checks if two string input fields are equal.
func (f *Form) IsEqualString(field1 string, field2 string) {
	string1 := f.Get(field1)
	string2 := f.Get(field2)
	if string1 != string2 {
		f.Errors.Add(field2, fmt.Sprintf("%s must be equal to %s", field2, field1))
	}
}

var EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// http://www.inanzzz.com/index.php/post/8l1a/validating-user-password-in-golang-requests
// SecurePassword defines a set of requirements for a password.
func (f *Form) SecurePassword(field string) {
	password := f.Get(field)
	var (
		upp, low, num, sym bool
		tot                uint8
	)

	for _, char := range password {
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
		}
	}

	if !upp || !low || !num || !sym || tot < 8 {
		f.Errors.Add(field, "Password length must be at least 8 characters and must contain at least 1 uppercase character, 1 lowercase character, 1 number, and 1 special character")
	}
}

// Valid returns true if the form has no errors.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
