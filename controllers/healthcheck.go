package controllers

import (
	"fmt"
	"kwoc20-backend/utils"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pong!")
	w.WriteHeader(200)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()

	if db != nil {
		fmt.Fprintf(w, "The server is up and database is reachable")
		w.WriteHeader(200)
	} else {
		fmt.Fprintf(w, "The database is unreachable")
		w.WriteHeader(200)
	}

	//Alternatively, code from utils.GetDB() could be replicated to give a more exact error message without the system panicking
}
