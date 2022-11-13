package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"kwoc20-backend/routes"
	"kwoc20-backend/utils"
)

func main() {
	// Set up logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	utils.InitialMigration()

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	router := mux.NewRouter()
	// if m := os.Getenv("MODE"); m == "dev" {
	// 	testSubRoute := router.PathPrefix("/test").Subrouter()
	// 	routes.RegisterTest(testSubRoute)
	// }

	// register logger middleware
	router.Use(utils.Logger)


	oauthSubRoute := router.PathPrefix("/oauth").Subrouter()
	routes.RegisterOAuth(oauthSubRoute)

	mentorSubRoute := router.PathPrefix("/mentor").Subrouter()
	routes.RegisterMentor(mentorSubRoute)

	studentSubRoute := router.PathPrefix("/student").Subrouter()
	routes.RegisterStudent(studentSubRoute)

	projectSubRoute := router.PathPrefix("/project").Subrouter()
	routes.RegisterProject(projectSubRoute)

	log.Info().Msg("Starting server on port " + port)

	router.PathPrefix("/").HandlerFunc(utils.PermissiveCORS).Methods("OPTIONS")

	router.MethodNotAllowedHandler = utils.GetMethodNotAllowedHandler()
	router.NotFoundHandler = utils.GetNotFoundHandler()

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal().Err(err).Msg("Error in starting server")
		os.Exit(1)
	}

}
