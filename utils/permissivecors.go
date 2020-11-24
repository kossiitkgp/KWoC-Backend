package utils

import "net/http"

//PermissiveCORS May not be needed in production (with same domain name for the frontend and backend)
func PermissiveCORS(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		next(w, r)
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
	}
}
