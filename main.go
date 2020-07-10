package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

	router.HandleFunc("/UserOAuth", controllers.MentorOAuth).Methods("POST")
	router.HandleFunc("/mentor", controllers.MentorReg).Methods("POST")
	router.HandleFunc("/project", controllers.ProjectReg).Methods("POST")
	router.HandleFunc("/project/all", controllers.ProjectGet).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+port,
		handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(router)))
}
