package routes

import (
	"kwoc20-backend/controllers"
	"kwoc20-backend/utils"

	"github.com/gorilla/mux"
)

func RegisterGetStats(r *mux.Router) {
	r.HandleFunc("/student/exists/{username}", utils.JsonIO(controllers.CheckStudent)).Methods("GET")

	r.HandleFunc("/students", utils.JsonIO(controllers.AllStudents)).Methods("GET")
	r.HandleFunc("/student/{username}", utils.LoginRequired(utils.JsonIO(controllers.OneStudent))).Methods("POST")

	r.HandleFunc("/projects", utils.JsonIO(controllers.GetAllProjects)).Methods("GET")
	r.HandleFunc("/mentor/{Mentor.Username}", utils.LoginRequired(utils.JsonIO(controllers.OneMentor))).Methods("POST")
}
