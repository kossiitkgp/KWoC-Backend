package controllers

import (
	"fmt"
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

	// Alternative code with better error reporting below, though this should probably be moved to utils

	/*
		DatabaseUsername := os.Getenv("DATABASE_USERNAME")
		DatabasePassword := os.Getenv("DATABASE_PASSWORD")
		DatabaseName := os.Getenv("DATABASE_NAME")
		DatabaseHost := os.Getenv("DATABASE_HOST")
		DatabasePort := os.Getenv("DATABASE_PORT")

		newURI := "host=" + DatabaseHost + " port=" + DatabasePort + " user=" + DatabaseUsername + " dbname=" + DatabaseName + " sslmode=disable password=" + DatabasePassword
		_, err := gorm.Open("postgres", newURI)
		if err != nil {
			fmt.Fprintf(w, "The database is ureachable with error: %s", err)
			w.WriteHeader(200)
		} else {
			fmt.Fprintf(w, "The server is up and database is reachable")
			w.WriteHeader(200)
		}
	*/
}
