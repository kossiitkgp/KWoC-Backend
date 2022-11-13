// Server utilities
package utils

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// Logging middleware for incoming requests
func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Info().Msgf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func getErrorHandler(errCode int, errMsg string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(errCode)

		log.Warn().Msgf(
			"Invalid Request: %s %s %d %s",
			r.Method,
			r.RequestURI,
			errCode,
			errMsg,
		)
	})
}

func GetNotFoundHandler() http.Handler {
	return getErrorHandler(404, "Path Not Found")
}

func GetMethodNotAllowedHandler() http.Handler {
	return getErrorHandler(405, "Method Not Allowed")
}