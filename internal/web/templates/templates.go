package templates

import (
	"html/template"
	"path/filepath"
	"strings"
)

// humanDate returns time in a friendlier format.
// date is a string in the format: "2021-02-24T13:35:50.028603852Z"
// Go sqlite cannot handle time.Time properly:
// https://github.com/mattn/go-sqlite3/issues/142
func humanDate(date string) string {
	if date == "" {
		return ""
	}
	t := strings.Split(date, "T")
	day := t[0]
	hours := strings.Split(t[1], ".")[0]
	return day + " at " + hours
}

var functions = template.FuncMap{"humanDate": humanDate}

// NewCache creates a new cache.
func NewCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts

	}

	return cache, nil
}
