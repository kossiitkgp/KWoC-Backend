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
		Name:     req["name"].(string),
		Email:    req["email"].(string),
		Username: req["username"].(string),
	}).Error
	if err != nil {
		return "Couldn't Connect to Database, Try Later", http.StatusInternalServerError
	}

	return "Successfully registered", http.StatusOK
}

func MentorDashboard(req map[string]interface{}, r *http.Request) (interface{}, int) {
	username := req["username"].(string)

	mentor := models.Mentor{}
	db := utils.GetDB()
	defer db.Close()

	db.Where(&models.Mentor{Username: username}).First(&mentor)

	if mentor.ID == 0 {
		return "User does not exist", http.StatusBadRequest
	}

	var projects []models.Project
	db.Where("mentor_id = ? OR secondary_mentor_id = ?", mentor.ID, mentor.ID).Preload("Mentor").Preload("SecondaryMentor").Find(&projects)

	// var secondary_projects []models.Project
	// db.Where("secondary_mentor_id = ?", mentor.ID).Find(&secondary_projects)

	all_projects := projects
	// projects_json, err := json.Marshal(projects)

	// var projects []models.Project
	// projects := models.Project{}
	// db.Where(&models.Project{MentorID: mentor.ID}).First(&projects)

	type Response map[string]interface{}
	res := Response{
		"name":     mentor.Name,
		"projects": all_projects,
	}

	return res, http.StatusOK
}

func GetAllMentors(req map[string]interface{}, r *http.Request) (interface{}, int) {
	db := utils.GetDB()
	defer db.Close()

	ctx_user := r.Context().Value(utils.CtxUserString("user")).(string)

	mentor := req["mentor"].(string)

	if ctx_user != mentor {
		utils.LOG.Printf("%v != %v Detected Session Hijacking\n", mentor, ctx_user)
		return "Corrupt JWT", http.StatusForbidden
	}

	var mentors []models.Mentor

	err := db.Select([]string{"Name", "Username"}).Not("username", mentor).Find(&mentors).Error
	if err != nil {
		utils.LOG.Println(err)
		return err, http.StatusInternalServerError
	}

	return mentors, http.StatusOK
}
