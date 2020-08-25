package routes

import (
	"github.com/gorilla/mux"

	"kwoc20-backend/controllers"
	"kwoc20-backend/utils"
)

func RegisterProject(r *mux.Router) {
	// Wrap the below Endpoint under LoginRequired Middleware after testing
	r.HandleFunc("/add", utils.LoginRequired(utils.JsonIO(controllers.ProjectReg))).Methods("POST")
	r.HandleFunc("/all", controllers.ProjectGet).Methods("GET")
}

