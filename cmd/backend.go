package main

import (
	"flag"
	"kwoc-backend/server"
	"kwoc-backend/utils"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Parse command-line arguments
	envFile := flag.String("envFile", ".env", "A file to load environment variables from.")
	flag.Parse()

	// Logger options ( using zerrolog )
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load environment variables via .env files
	log.Info().Msgf("Attempting to load environment variables from %s.", *envFile)
	dotenv_err := godotenv.Load(*envFile)

	if dotenv_err != nil {
		log.Warn().Msgf("Failed to load environment variables from %s.", *envFile)
	} else {
		log.Info().Msgf("Successfully loaded environment variables from %s.", *envFile)
	}

	mig_err := utils.MigrateModels()
	if mig_err != nil {
		log.Fatal().Err(mig_err).Msg("Database migration error.")
	}

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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	log.Info().Msg("Starting server on port : " + port)
	err = http.ListenAndServe(":"+port, router)

	if err != nil {
		log.Fatal().Err(err).Msg("Error starting the server.")
	}
}

func cleanup() {
	log.Info().Msg("Received SIGINT, Shutting down server")
}
