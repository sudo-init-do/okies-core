package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/sudo-init-do/okies_core/pkg/response"
)

// RecoveryMiddleware recovers from panics and returns JSON error
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// Return safe JSON response
				response.Write(w, http.StatusInternalServerError, "internal server error", nil)
				// Print stack trace for debugging
				debug.PrintStack()
			}
		}()
		next.ServeHTTP(w, r)
	})
}
