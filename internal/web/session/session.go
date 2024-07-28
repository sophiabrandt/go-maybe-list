package session

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

// New creates the session manager.
func New() *scs.SessionManager {
	session := scs.New()
	session.Lifetime = 12 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteStrictMode
	session.Cookie.Secure = true
	return session
}
