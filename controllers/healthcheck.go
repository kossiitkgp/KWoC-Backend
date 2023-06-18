package controllers

import (
	"kwoc-backend/utils"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// Ping responds with "pong" and returns the latency.
func Ping(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("pong"))
	if err != nil {
		log.Err(err).Msg("Could not respond to Ping")
	}

	elapsed := time.Since(start)
	log.Info().Str("latency", elapsed.String()).Msg("Ping request processed")
}

// HealthCheck checks the server and database status.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	db, err := utils.GetDB()
	if err != nil {
		log.Err(err).Msg("Could not connect to database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.Exec("SELECT 1").Error
	if err != nil {
		log.Err(err).Msg("Could not ping database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("OK"))
	if err != nil {
		log.Err(err).Msg("Could not respond to HealthCheck")
	}

	log.Info().Msg("Healthcheck request is OK")
}
