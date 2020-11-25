package controllers

import (
	"net/http"

	"kwoc20-backend/models"
	"kwoc20-backend/utils"
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
	username := req["username"].(string)

	
	mentor := models.Mentor{}
	db := utils.GetDB()
	defer db.Close()
	db.Where(&models.Mentor{Username: username}).First(&mentor)
	if mentor.ID == 0 {
		return "no user", 400
	}
	
	var projects []models.Project
	db.Where("mentor_id = ?", mentor.ID).Find(&projects)
	// projects_json, err := json.Marshal(projects)
	
	// var projects []models.Project
	// projects := models.Project{}
	// db.Where(&models.Project{MentorID: mentor.ID}).First(&projects)

	type Response map[string]interface{}
	res := Response{
		"name": mentor.Name,
		"projects": projects,
	}

	return res, 200
}

