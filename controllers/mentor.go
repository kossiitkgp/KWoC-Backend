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