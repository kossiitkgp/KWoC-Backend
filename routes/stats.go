package routes

import (
	"kwoc20-backend/controllers"
	"kwoc20-backend/utils"

	"github.com/gorilla/mux"
)

func RegisterGetStats(r *mux.Router) {
	r.HandleFunc("/student/exists/{username}", utils.PermissiveCORSMiddleware(utils.JsonIO(controllers.CheckStudent))).Methods("GET")

	r.HandleFunc("/students", utils.PermissiveCORSMiddleware(utils.JsonIO(controllers.AllStudents))).Methods("GET")
	r.HandleFunc("/student/{username}", utils.PermissiveCORSMiddleware(utils.LoginRequired(utils.JsonIO(controllers.OneStudent)))).Methods("POST")

	r.HandleFunc("/projects", utils.PermissiveCORSMiddleware(utils.JsonIO(controllers.GetAllProjects))).Methods("GET")
	r.HandleFunc("/mentor/{Mentor.Username}", utils.PermissiveCORSMiddleware(utils.LoginRequired(utils.JsonIO(controllers.OneMentor)))).Methods("POST")

	r.HandleFunc("/overall", utils.PermissiveCORSMiddleware(utils.JsonIO(controllers.OverallStats))).Methods("GET")
}
