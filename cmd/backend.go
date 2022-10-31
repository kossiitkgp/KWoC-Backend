package main

import (
	"kwoc-backend/server"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Logger options ( using zerrolog )
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("Creating mux router")
	router := server.NewRouter()

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	// Sanity check for backend port
	_, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal().Err(err).Msg("BACKEND_PORT env variable is invalid")
		os.Exit(1)
	}

	// Handling INTERRUPT signal for cleanup in a new goroutine.
	// This is not necessary, but good for log keeping
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	log.Info().Msg("Starting server on port : " + port)
	http.ListenAndServe(":"+port, router)
}

func cleanup() {
	log.Info().Msg("Received SIGINT, Shutting down server")
}
