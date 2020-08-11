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

