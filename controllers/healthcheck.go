package controllers

import (
	"net/http"

	"kwoc20-backend/utils"

	"github.com/rs/zerolog/log"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, err := w.Write([]byte("Pong!"))
	if err != nil {
		log.Info().Msg("Could not respond to Ping")
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	w.WriteHeader(200)
	if db != nil {
		_, err := w.Write([]byte("The server is up and database is reachable"))
		if err != nil {
			log.Info().Msg("Could not respond to Ping")
		}
		return
	}
	_, err := w.Write([]byte("Database is unreachable"))
	if err != nil {
		log.Info().Msg("Could not respond to Ping")
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
