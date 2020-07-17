package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	logs "kwoc20-backend/utils/logs/pkg"

	"github.com/go-kit/kit/log/level"
	"github.com/jinzhu/gorm"

	"kwoc20-backend/models"
)

//ProjectReg endpoint to register project details
func ProjectReg(w http.ResponseWriter, r *http.Request) {

	var project models.Project
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &project)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		level.Error(logs.Logger).Log("error", fmt.Sprintf("%v", err))
		return
	}

	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		level.Error(logs.Logger).Log("error", fmt.Sprintf("%v", err))
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		level.Error(logs.Logger).Log("error", fmt.Sprintf("%v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "success"}`))

}

//ProjectGet endpoint to fetch all projects
// INCOMPLETE BECAUSE MENTOR STILL NEEDS TO BE ADDED
func ProjectGet(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		level.Error(logs.Logger).Log("error", fmt.Sprintf("%v", err))
		return
	}
	defer db.Close()

	var projects []models.Project
	err = db.Find(&projects).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		level.Error(logs.Logger).Log("error", fmt.Sprintf("%v", err))
		return
	}

	err = json.NewEncoder(w).Encode(projects)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		level.Error(logs.Logger).Log("error", fmt.Sprintf("%v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "message" : "success" }`))

}
