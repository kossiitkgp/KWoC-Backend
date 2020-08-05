package routes

import (
	"github.com/gorilla/mux"

	"kwoc20-backend/controllers"
)

func RegisterMentor(r *mux.Router) {
	r.HandleFunc("/", controllers.MentorReg).Methods("POST")
}

