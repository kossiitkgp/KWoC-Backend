package controllers

import (
	"kwoc20-backend/controllers"

	"github.com/gorilla/mux"
)

func RegisterGetStats(r *mux.Router) {

	r.HandleFunc("/stats/student/{Username}", CheckStudent).Methods("GET")
	r.HandleFunc("/stats/students", controllers.AllStudents).Methods("GET")
	r.HandleFunc("/stats/student/{username}", controllers.OneStudent).Methods("GET")
	r.HandleFunc("/stats/projects", controllers.AllProjects).Methods("GET")
	r.HandleFunc("/stats/projects/{Mentor.Username}", OneMentor).Methods("GET")
}
