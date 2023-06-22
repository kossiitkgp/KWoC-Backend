package controllers

import (
	"kwoc-backend/middleware"
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
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	err := db.Exec("SELECT 1").Error
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
