package routes

import (
	"github.com/gorilla/mux"

	"kwoc20-backend/controllers"
)

func RegisterProject(r *mux.Router) {
	r.HandleFunc("/", controllers.ProjectReg).Methods("POST")
	r.HandleFunc("/all", controllers.ProjectGet).Methods("GET")
}

