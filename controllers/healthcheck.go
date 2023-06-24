package controllers

import (
	"fmt"
	"kwoc-backend/middleware"
	"kwoc-backend/utils"
	"net/http"
	"time"
)

// Ping responds with "pong" and returns the latency.
func Ping(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("pong"))
	if err != nil {
		utils.LogErr(r, err, "Could not respond to Ping")
		return
	}

	elapsed := fmt.Sprint(time.Since(start))
	var Info string = "latency: " + elapsed + " Ping request processed"
	utils.LogInfo(r, Info)
}

// HealthCheck checks the server and database status.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	err := db.Exec("SELECT 1").Error
	if err != nil {
		utils.LogErr(r, err, "Could not ping database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("OK"))
	if err != nil {
		utils.LogErr(r, err, "Could not respond to HealthCheck")
		return
	}

	utils.LogInfo(r, "Healthcheck request is OK")
}
