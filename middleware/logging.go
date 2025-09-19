package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs requests with method, path, status and duration
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("➡️ %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("✅ %s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	})
}
