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

func GetNotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)

		log.Warn().Msgf(
			"Invalid Request: %s %s %s",
			r.Method,
			r.RequestURI,
			"404 Path Not Found",
		)
	})
}

func GetMethodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)

		log.Warn().Msgf(
			"Invalid Request: %s %s %s",
			r.Method,
			r.RequestURI,
			"405 Method Not Allowed",
		)
	})
}