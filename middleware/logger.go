package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kossiitkgp/kwoc-backend/utils"
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
