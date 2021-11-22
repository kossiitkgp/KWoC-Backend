package routes

import (
	"github.com/gorilla/mux"

	"kwoc20-backend/controllers"
	"kwoc20-backend/utils"
)

func RegisterProject(r *mux.Router) {
	// Wrap the below Endpoint under LoginRequired Middleware after testing
	r.HandleFunc("/add", utils.PermissiveCORSMiddleware(utils.LoginRequired(utils.JsonIO(controllers.ProjectReg)))).Methods("POST")
	r.HandleFunc("/all", utils.PermissiveCORSMiddleware(controllers.AllProjects)).Methods("GET")
	r.HandleFunc("/stats", utils.PermissiveCORSMiddleware(utils.JsonIO(controllers.RunStats))).Methods("GET")
	r.HandleFunc("/get", utils.PermissiveCORSMiddleware(utils.LoginRequired(utils.JsonIO(controllers.ProjectDetails)))).Methods("GET")
	r.HandleFunc("/update", utils.PermissiveCORSMiddleware(utils.LoginRequired(utils.JsonIO(controllers.UpdateDetails)))).Methods("PUT")
}
