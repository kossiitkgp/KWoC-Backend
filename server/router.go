// Server startup and routes
package server

import (
	"net/http"

	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Setup up mux routes and router
func NewRouter(db *gorm.DB) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.LogWarn(r, "404 Not Found.")
	})

	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.LogWarn(r, "405 Method Not Allowed.")
	})

	// iterate over all routes
	app := &middleware.App{Db: db}
	routes := getRoutes(app)

	for _, route := range routes {
		var handler http.Handler

		// logger middleware to log incoming requests
		handler = route.HandlerFunc
		handler = middleware.Logger(handler, route.Name)

		// register route
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
