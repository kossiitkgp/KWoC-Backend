package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"kwoc20-backend/routes"
	"kwoc20-backend/utils"
)

func main() {

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

	oauthSubRoute := router.PathPrefix("/oauth").Subrouter()
	routes.RegisterOAuth(oauthSubRoute)

	mentorSubRoute := router.PathPrefix("/mentor").Subrouter()
	routes.RegisterMentor(mentorSubRoute)

	studentSubRoute := router.PathPrefix("/student").Subrouter()
	routes.RegisterStudent(studentSubRoute)

	projectSubRoute := router.PathPrefix("/project").Subrouter()
	routes.RegisterProject(projectSubRoute)

	var mainLogger = log.New(os.Stderr, "Message: ", log.LstdFlags|log.Lshortfile)
	mainLogger.Println("Starting server on port " + port)

	err := http.ListenAndServe(":"+port,
		router)
	if err != nil {
		mainLogger.Println("Error in Starting ", err)
		os.Exit(1)
	}

}
