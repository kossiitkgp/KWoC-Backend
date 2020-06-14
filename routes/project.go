package routes

import (
	"net/http"
	"encoding/json"
	"io/ioutil"

	"github.com/jinzhu/gorm"

	"kwoc20-backend/models"
)

// endpoint to register project details
func ProjectReg(w http.ResponseWriter, r *http.Request) {
	
	var project models.Project
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &project)
	if err != nil {
		panic(err)
	}

	
	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		panic("failed to connect DB")
	}
	defer db.Close()

	// need to add error handling here
	db.Create(&models.Project{
		Name: project.Name,
		Desc: project.Desc,
		Tags: project.Tags,
		RepoLink: project.RepoLink,
		ComChannel: project.ComChannel,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "success"}`))
}


// endpoint to fetch all projects
// INCOMPLETE BECAUSE MENTOR STILL NEEDS TO BE ADDED
func ProjectGet(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "kwoc.db")
    if err != nil {
        panic("failed to connect database")
    }
    defer db.Close()

    var projects []models.Project
    db.Find(&projects)
	
	w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(projects)
}

