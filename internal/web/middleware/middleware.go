// https://blog.questionable.services/article/guide-logging-middleware-go/
// https://lets-go.alexedwards.net/
package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/pkg/errors"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

// LogRequest logs information about each request.
func LogRequest(log *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Printf(
						"ERROR: %s, trace: %s", err, debug.Stack(),
					)
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			log.Printf(
				"status: %d, method: %s, path: %s, duration: %s", wrapped.status, r.Method, r.URL.EscapedPath(), time.Since(start),
			)
		}

		return http.HandlerFunc(fn)
	}
}

// RecoverPanic closes a connection and returns an error response.
func RecoverPanic(log *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					err := errors.Errorf("%v", r)
					log.Printf(
						"PANIC: %s, trace: %s", err, debug.Stack(),
					)
					// Set a "Connection: close" header on the response.
					w.Header().Set("Connection", "close")
					// return internal server error
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// SecureHeaders sets header options.
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		next.ServeHTTP(w, r)
	})
}
