package main

import (
	"fmt"
	"net/http"
	"os"
	"log"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	
	"kwoc20-backend/routes"
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
	
	router.HandleFunc("/mentor", routes.MentorReg).Methods("POST")
	router.HandleFunc("/project", routes.ProjectReg).Methods("POST")
	router.HandleFunc("/project/all", routes.ProjectGet).Methods("GET")
	
	log.Fatal(http.ListenAndServe(":"+port, router))
}

