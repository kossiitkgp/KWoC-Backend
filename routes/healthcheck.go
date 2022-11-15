package routes

import (
	"github.com/gorilla/mux"
	"kwoc20-backend/controllers"
)

func HealthCheck(r *mux.Router) {
	// Wrap the below Endpoint under LoginRequired Middleware after testing
	r.HandleFunc("/ping", controllers.Ping).Methods("GET")
}
