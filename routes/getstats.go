package routes

import (
	"kwoc20-backend/controllers"

	"github.com/gorilla/mux"
)

func RegisterGetStats(r *mux.Router) {

	r.HandleFunc("/student/<username>", controllers.CheckStudent).Methods("GET")
	r.HandleFunc("/students", controllers.AllStudents).Methods("GET")
	r.HandleFunc("/student/<username>", controllers.OneStudent).Methods("GET")
	r.HandleFunc("/projects", controllers.GetAllProjects).Methods("GET")
	r.HandleFunc("/projects/<Mentor.Username>", controllers.OneMentor).Methods("GET")
}
