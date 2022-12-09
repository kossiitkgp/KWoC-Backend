package routes

import (
	"kwoc20-backend/controllers"
	"kwoc20-backend/utils"

	"github.com/gorilla/mux"
)

func RegisterGetStats(r *mux.Router) {

	r.HandleFunc("/student/exists/{username}", utils.JsonIO(controllers.CheckStudent)).Methods("GET")
	r.HandleFunc("/students", controllers.AllStudents).Methods("GET")
	r.HandleFunc("/student/{username}", controllers.OneStudent).Methods("GET")
	r.HandleFunc("/projects", controllers.GetAllProjects).Methods("GET")
	r.HandleFunc("/projects/{Mentor.Username}", controllers.OneMentor).Methods("GET")
}
