package controllers

import (
	"net/http"

	"kwoc20-backend/models"
	"kwoc20-backend/utils"
	"encoding/json"
	"github.com/gorilla/mux"
)

// After Being checked by LoginRequired Middleware
func MentorReg(req map[string]interface{}, r *http.Request) (interface{}, int) {
	db := utils.GetDB()
	defer db.Close()

	err := db.Create(&models.Mentor{
		Name:         req["name"].(string),
		Email:        req["email"].(string),
		Username: 	  req["username"].(string),
	}).Error

	if err != nil {
		return "database issue", 500
	}

	return "success", 200
}

func MentorDashboard(req map[string]interface{}, r *http.Request) (interface{}, int){
	// return "name", 200
	username := req["username"].(string)

	mentor := models.Mentor{}
	db := utils.GetDB()
	defer db.Close()
	db.Where("name = ?", username).First(&mentor)

	if mentor.ID == 0 {
		return "no user", 400
	}

	type Response map[string]interface{}
	res := Response{
		"username" : username,
		"name": mentor.Name,
	}

	return res, 200
}

func MentorProjects(w http.ResponseWriter, r *http.Request)() {

	db := utils.GetDB()
	defer db.Close()

	vars := mux.Vars(r)	
	mentor_username := vars["MENTOR_USERNAME"]

	var mentor models.Mentor 
	db.Where("username = ?", mentor_username).First(&mentor)
	mentor_id := mentor.ID

	if mentor_id == 0 {
		w.WriteHeader(500)
		return
	}
	var projects []models.Project
	db.Where("mentor_id = ?", mentor_id).Find(&projects)
	projects_json, err := json.Marshal(projects)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(projects_json)
}