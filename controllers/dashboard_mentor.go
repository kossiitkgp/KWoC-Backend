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

type ProjectInfo struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Tags          string `json:"tags"`
	RepoLink      string `json:"repo_link"`
	CommChannel   string `json:"comm_channel"`
	ReadmeLink    string `json:"readme_link"`
	ProjectStatus bool   `json:"project_status"`

	// stats table
	CommitCount  uint `json:"commit_count"`
	PullCount    uint `json:"pull_count"`
	LinesAdded   uint `json:"lines_added"`
	LinesRemoved uint `json:"lines_removed"`
}

type MentorDashboard struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`

	Projects []ProjectInfo      `json:"projects"`
	Students []StudentDashboard `json:"students"`
}

func CreateMentorDashboard(mentor models.Mentor, db *gorm.DB) MentorDashboard {
	var projects []models.Project
	var projectsInfo []ProjectInfo
	var students []StudentDashboard

	db.Table("projects").
		Where("mentor_id = ? OR secondary_mentor_id = ?", mentor.ID, mentor.ID).
		Preload("Mentor").Preload("SecondaryMentor").
		Find(&projects)

	for _, project := range projects {
		projectInfo := ProjectInfo{
			Name:          project.Name,
			Description:   project.Description,
			Tags:          project.Tags,
			RepoLink:      project.RepoLink,
			CommChannel:   project.CommChannel,
			ReadmeLink:    project.ReadmeLink,
			ProjectStatus: project.ProjectStatus,

			// stats table
			CommitCount:  project.CommitCount,
			PullCount:    project.PullCount,
			LinesAdded:   project.LinesAdded,
			LinesRemoved: project.LinesRemoved,
		}
		projectsInfo = append(projectsInfo, projectInfo)

		var modelStudent models.Student
		for _, studentUsername := range strings.Split(project.Contributors, ",") {
			db.Table("students").Where("username = ?", studentUsername).First(&modelStudent)
			student := CreateStudentDashboard(modelStudent, db)
			students = append(students, student)
		}
	}

	return MentorDashboard{
		Name:     mentor.Name,
		Username: mentor.Username,
		Email:    mentor.Email,

		Projects: projectsInfo,
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
