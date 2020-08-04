package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterTest(r *mux.Router) {
	r.HandleFunc("/oauth", func(w http.ResponseWriter, r *http.Request) {
		path := "./tests/oauth.html"
		http.ServeFile(w, r, path)
	})
}
