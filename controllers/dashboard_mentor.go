package controllers

import (
	"fmt"
	"kwoc-backend/middleware"
	"kwoc-backend/models"
	"kwoc-backend/utils"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type MentorDashboard struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`

	Projects []models.Project `json:"projects"`
	Students []models.Student `json:"students"`
}

func CreateMentorDashboard(mentor models.Mentor, db *gorm.DB) MentorDashboard {
	var projects []models.Project
	var students []models.Student

	db.Table("projects").
		Where("mentor_id = ? OR secondary_mentor_id = ?", mentor.ID, mentor.ID).
		Preload("Mentor").Preload("SecondaryMentor").
		Find(&projects)

	for _, project := range projects {
		var student models.Student
		for _, studentUsername := range strings.Split(project.Contributors, ",") {
			db.Table("students").Where("username = ?", studentUsername).First(&student)
			students = append(students, student)
		}
	}

	return MentorDashboard{
		Name:     mentor.Name,
		Username: mentor.Username,
		Email:    mentor.Email,

		Projects: projects,
		Students: students,
	}
}

func FetchMentorDashboard(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	var modelMentor models.Mentor

	login_username := r.Context().Value(middleware.LoginCtxKey(middleware.LOGIN_CTX_USERNAME_KEY))
	tx := db.
		Table("mentors").
		Where("username = ?", login_username).
		First(&modelMentor)

	if tx.Error == gorm.ErrRecordNotFound {
		utils.LogErrAndRespond(
			r,
			w,
			tx.Error,
			fmt.Sprintf("Mentor `%s` does not exists.", login_username),
			http.StatusBadRequest,
		)
		return
	}
	if tx.Error != nil {
		utils.LogErrAndRespond(
			r,
			w,
			tx.Error,
			fmt.Sprintf("Database Error fetching mentor with username `%s`", login_username),
			http.StatusInternalServerError,
		)
		return
	}

	mentor := CreateMentorDashboard(modelMentor, db)

	utils.RespondWithJson(r, w, mentor)
}
