package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"

	"kwoc20-backend/models"
	utils "kwoc20-backend/utils"
	
)

//ProjectReg endpoint to register project details
func ProjectReg(w http.ResponseWriter, r *http.Request) {

	var project models.Project
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &project)
	if err != nil {
		http.Error(w, err.Error(), 400)
		utils.LOG.Println(err)
		return
	}

	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		http.Error(w, err.Error(), 500)
		utils.LOG.Println(err)
		return
	}
	defer db.Close()

	err = db.Create(&models.Project{
		Name:       project.Name,
		Desc:       project.Desc,
		Tags:       project.Tags,
		RepoLink:   project.RepoLink,
		ComChannel: project.ComChannel,
	}).Error

	if err != nil {
		utils.LOG.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`success`))

}

//ProjectGet endpoint to fetch all projects
// INCOMPLETE BECAUSE MENTOR STILL NEEDS TO BE ADDED
func ProjectGet(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		http.Error(w, err.Error(), 500)
		utils.LOG.Println(err)
		return
	}
	defer db.Close()

	var projects []models.Project
	err = db.Find(&projects).Error
	if err != nil {
		http.Error(w, err.Error(), 400)
		utils.LOG.Println(err)
		return
	}

	err = json.NewEncoder(w).Encode(projects)
	if err != nil {
		http.Error(w, err.Error(), 500)
		utils.LOG.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`success`))

}
