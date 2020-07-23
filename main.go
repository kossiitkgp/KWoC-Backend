package main

import (
	"fmt"
	"net/http"
	"os"
	"log"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"kwoc20-backend/controllers"
	"kwoc20-backend/models"
)

func initialMigration() {
	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&models.Project{})
}

func main() {

	initialMigration()

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/oauth", controllers.UserOAuth).Methods("POST")
	router.HandleFunc("/mentor", controllers.MentorReg).Methods("POST")
	router.HandleFunc("/project", controllers.ProjectReg).Methods("POST")
	router.HandleFunc("/project/all", controllers.ProjectGet).Methods("GET")

	var mainLogger = log.New(os.Stderr, "Message: ", log.LstdFlags | log.Lshortfile)
	mainLogger.Println("Starting server on port "+port)
	
	err := http.ListenAndServe(":"+port,
		handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(router))
	if err != nil {
		mainLogger.Println("Error in Starting ",err)
		os.Exit(1)
	}

	
}
