package routes

import (
	"kwoc20-backend/controllers"

	"github.com/gorilla/mux"
)

func RegisterGetStats(r *mux.Router) {

	r.HandleFunc("/stats/student/{Username}", controllers.CheckStudent).Methods("GET")
	r.HandleFunc("/stats/students", controllers.AllStudents).Methods("GET")
	r.HandleFunc("/stats/student/{username}", controllers.OneStudent).Methods("GET")
	r.HandleFunc("/stats/projects", controllers.GetAllProjects).Methods("GET")
	r.HandleFunc("/stats/projects/{Mentor.Username}", controllers.OneMentor).Methods("GET")
}
