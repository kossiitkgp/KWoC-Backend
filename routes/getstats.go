package routes

import (
	"kwoc20-backend/controllers"

	"github.com/gorilla/mux"
)

func RegisterHealthCheck(r *mux.Router) {

	r.HandleFunc("/stats/students", controllers.AllStudents).Methods("GET")
	r.HandleFunc("/stats/student/{username}", controllers.OneStudent).Methods("GET")
	r.HandleFunc("/stats/projects", controllers.AllProjects).Methods("GET")
}
