// Server utilities
package utils

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// Logging middleware for incoming requests
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Info().Msgf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
