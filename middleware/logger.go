package middleware

import (
	"fmt"
	"kwoc-backend/utils"
	"net/http"
	"time"
)

// Logging middleware for incoming requests
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		utils.LogInfo(
			r,
			fmt.Sprintf("%s %s", name, time.Since(start)),
		)
	})
}
