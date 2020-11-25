package controllers

import (
	"fmt"
	"net/http"

	"kwoc20-backend/models"
	"kwoc20-backend/utils"
)

// After Being checked by LoginRequired Middleware
func StudentReg(req map[string]interface{}, r *http.Request) (interface{}, int) {
	db := utils.GetDB()
	defer db.Close()

	err := db.Create(&models.Student{
		Name:         req["name"].(string),
		Email:        req["email"].(string),
		College:      req["college"].(string),
		Username: 	  req["username"].(string),
	}).Error

	if err != nil {
		fmt.Println("err is ",err)
		return "database issue", 500
	}

	return "success", 200
}

func StudentDashboard(req map[string]interface{}, r *http.Request) (interface{}, int){
	// return "name", 200
	username := req["username"].(string)

	student := models.Student{}
	db := utils.GetDB()
	defer db.Close()
	db.Where(&models.Student{Username: username}).First(&student)
	if student.ID == 0 {
		return "no user", 400
	}

	type Response map[string]interface{}
	res := Response{
		"name": student.Name,
		"college": student.College,
	}

	return res, 200
	
}

func StudentStats (username string) interface{} {
	return fmt.Sprintf("stats of %s", username)
}