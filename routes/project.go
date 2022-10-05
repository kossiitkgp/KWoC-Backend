package routes

import (
	"kwoc20-backend/controllers"
	"kwoc20-backend/utils"

	"github.com/gorilla/mux"
)

func RegisterProject(r *mux.Router) {
	// Wrap the below Endpoint under LoginRequired Middleware after testing
	r.HandleFunc("", utils.PermissiveCORSMiddleware(utils.JsonIO(controllers.AllProjects))).Methods("GET")
	r.HandleFunc("/register", utils.PermissiveCORSMiddleware(utils.LoginRequired(utils.JsonIO(controllers.ProjectReg)))).Methods("POST")
	r.HandleFunc("/stats", utils.PermissiveCORSMiddleware(utils.JsonIO(controllers.RunStats))).Methods("GET")
	r.HandleFunc("/details", utils.PermissiveCORSMiddleware(utils.LoginRequired(utils.JsonIO(controllers.ProjectDetails)))).Methods("POST")
	r.HandleFunc("/update", utils.PermissiveCORSMiddleware(utils.LoginRequired(utils.JsonIO(controllers.UpdateDetails)))).Methods("PUT")
}
