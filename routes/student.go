package routes

import (
	"github.com/gorilla/mux"

	// TEMP
	"kwoc20-backend/controllers"
	"kwoc20-backend/utils"
)

// TEMP
// Discuss and add 2 middlewares - JWT Required, JSON Marshalling
func RegisterStudent(r *mux.Router) {
	r.HandleFunc("/form", utils.PermissiveCORSMiddleware(utils.LoginRequired(utils.JsonIO(controllers.StudentReg)))).Methods("POST")
	r.HandleFunc("/dashboard", utils.PermissiveCORSMiddleware(utils.JsonIO(controllers.StudentDashboard))).Methods("POST")
	r.HandleFunc("/blog", utils.PermissiveCORSMiddleware(utils.LoginRequired(utils.JsonIO(controllers.StudentBlogLink)))).Methods("POST")
	// r.HandleFunc("/form", utils.JsonIO(controllers.StudentReg)).Methods("POST")
}
