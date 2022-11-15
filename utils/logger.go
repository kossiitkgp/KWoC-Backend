// Server utilities
package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func LogErr(r *http.Request, err error, errMsg string) {
	log.Err(err).Msgf(
		"%s %s: %s",
		r.Method,
		r.RequestURI,
		errMsg,
	)
}

func LogInfo(r *http.Request, info string) {
	log.Info().Msgf(
		"%s %s: %s",
		r.Method,
		r.RequestURI,
		info,
	)
}

func LogWarn(r *http.Request, warning string) {
	log.Warn().Msgf(
		"%s %s: %s",
		r.Method,
		r.RequestURI,
		warning,
	)
}

// Logging middleware for incoming requests
func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		LogInfo(
			r,
			fmt.Sprintf("Handled in %v", time.Since(start)),
		)
	})
}

func getErrorHandler(errCode int, errMsg string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(errCode)

		LogWarn(
			r,
			fmt.Sprintf("Invalid Request %d %s", errCode, errMsg),
		)
	})
}

func GetNotFoundHandler() http.Handler {
	return getErrorHandler(404, "Path Not Found")
}

func GetMethodNotAllowedHandler() http.Handler {
	return getErrorHandler(405, "Method Not Allowed")
}