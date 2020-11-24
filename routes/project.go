package routes

import (
	"github.com/gorilla/mux"

	"kwoc20-backend/controllers"
	"kwoc20-backend/utils"
)

func RegisterProject(r *mux.Router) {
	// Wrap the below Endpoint under LoginRequired Middleware after testing
	r.HandleFunc("/add", utils.PermissiveCORSMiddleware(utils.LoginRequired(utils.JsonIO(controllers.ProjectReg)))).Methods("POST")
	r.HandleFunc("/all", utils.PermissiveCORSMiddleware(controllers.ProjectGet)).Methods("GET")
}
