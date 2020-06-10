package main

import (
	"net/http"
	"os"
	"log"

	"kwoc20-backend/routes"
)



func main() {
	
	port := os.Getenv("PORT")
    if port == "" {
        port = "5000"
	}
	
	http.HandleFunc("/", routes.MentorReg)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

