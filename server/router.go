// Server startup and routes
package server

import (
	"kwoc-backend/utils"
	"net/http"

	"github.com/gorilla/mux"
)

// Setup up mux routes and router
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

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
