// Server startup and routes
package server

import (
	"errors"
	"net/http"

	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var ErrRouteNotFound = errors.New("route not found")
var ErrMethodNotAllowed = errors.New("method not allowed")

// Setup up mux routes and router
func NewRouter(db *gorm.DB, testMode bool) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.LogErrAndRespond(r, w, ErrRouteNotFound, "404 Not Found.", http.StatusNotFound)

	})

	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.LogErrAndRespond(r, w, ErrMethodNotAllowed, "405 Method Not Allowed.", http.StatusMethodNotAllowed)
	})

	// iterate over all routes
	app := &middleware.App{Db: db}
	routes := getRoutes(app)

	for _, route := range routes {
		// skip disabled routes
		if !testMode && route.disabled {
			continue
		}

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
