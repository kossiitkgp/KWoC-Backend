// Server startup and routes
package server

import (
	"kwoc-backend/utils"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Setup up mux routes and router
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Warn().Msgf(
			"%s %s Not Found",
			r.Method,
			r.RequestURI,
		)
	})

	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Warn().Msgf(
			"%s %s Method Not Allowed",
			r.Method,
			r.RequestURI,
		)
	})

	// iterate over rall routes
	for _, route := range routes {
		var handler http.Handler

		// logger middleware to log incoming requests
		handler = route.HandlerFunc
		handler = utils.Logger(handler, route.Name)

		// register route
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
