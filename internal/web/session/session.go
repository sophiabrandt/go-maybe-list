package session

import (
	"net/http"
	"time"

	"github.com/golangcollege/sessions"
)

// New creates the session manager.
func New(secret string) *sessions.Session {
	session := sessions.New([]byte(secret))
	session.Lifetime = 12 * time.Hour
	session.Persist = true
	session.SameSite = http.SameSiteStrictMode
	return session
}
