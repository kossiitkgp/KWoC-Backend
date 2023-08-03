package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"
)

// Ping godoc
//
//	@Summary		ping
//	@Description	Ping responds with "pong" and returns the latency
//	@Accept			plain
//	@Produce		plain
//	@Success		200	{string}	string	"pong"
//	@Router			/healthcheck/ping [get]
func Ping(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("pong"))
	if err != nil {
		utils.LogErr(r, err, "Could not respond to Ping")
		return
	}

	elapsed := time.Since(start)
	utils.LogInfo(r, fmt.Sprintf("latency: %dns Ping request processed", elapsed))
}

// HealthCheck godoc
//
//	@Summary		Checks the health status of the server and database.
//	@Description	The HealthCheck endpoint examines the operational status of the server and the associated database.
//	@Accept			plain
//	@Produce		plain
//	@Success		200	{string}	string	"OK"
//	@Failure		500	{string}	string	"Could not ping the database"
//	@Router			/healthcheck/ [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	err := db.Exec("SELECT 1").Error
	if err != nil {
		utils.LogErrAndRespond(r, w, err, "Could not ping database", http.StatusInternalServerError)
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
