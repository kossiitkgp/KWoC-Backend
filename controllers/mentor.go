package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"

	"github.com/kossiitkgp/kwoc-backend/v2/models"

	"gorm.io/gorm"
)

type RegisterMentorReqFields struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type ProjectInfo struct {
	Name          string `json:"name"`
	RepoLink      string `json:"repo_link"`
	ProjectStatus bool   `json:"project_status"`

	CommitCount  uint `json:"commit_count"`
	PullCount    uint `json:"pull_count"`
	LinesAdded   uint `json:"lines_added"`
	LinesRemoved uint `json:"lines_removed"`

	Pulls []string `json:"pulls"`
}

type StudentInfo struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

type MentorDashboard struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`

	Projects []ProjectInfo `json:"projects"`
	Students []StudentInfo `json:"students"`
}

func RegisterMentor(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db
	var reqFields = RegisterMentorReqFields{}

	err := utils.DecodeJSONBody(r, &reqFields)
	if err != nil {
		utils.LogErrAndRespond(r, w, err, "Error decoding JSON body.", http.StatusBadRequest)
		return
	}

	// Check if the JWT login username is the same as the mentor's given username
	login_username := r.Context().Value(middleware.LOGIN_CTX_USERNAME_KEY).(string)

	if reqFields.Username != login_username {
		utils.LogWarn(
			r,
			fmt.Sprintf(
				"POSSIBLE SESSION HIJACKING\nJWT Username: %s, Given Username: %s",
				login_username,
				reqFields.Username,
			),
		)

		utils.RespondWithHTTPMessage(r, w, http.StatusUnauthorized, "Login username and given username do not match.")
		return
	}

	// Check if the mentor already exists in the db
	mentor := models.Mentor{}
	tx := db.
		Table("mentors").
		Where("username = ?", reqFields.Username).
		First(&mentor)

	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		utils.LogErrAndRespond(r, w, tx.Error, "Database error.", http.StatusInternalServerError)
		return
	}

	mentor_exists := mentor.Username == reqFields.Username

	if mentor_exists {
		utils.LogWarnAndRespond(
			r,
			w,
			fmt.Sprintf("Mentor `%s` already exists.", mentor.Username),
			http.StatusBadRequest,
		)

		return
	}

	// Create a db entry if the mentor doesn't exist
	tx = db.Create(&models.Mentor{
		Username: reqFields.Username,
		Name:     reqFields.Name,
		Email:    reqFields.Email,
	})

	if tx.Error != nil {
		utils.LogErrAndRespond(r, w, tx.Error, "Database error.", http.StatusInternalServerError)
		return
	}

	utils.RespondWithHTTPMessage(r, w, http.StatusOK, "Mentor registration successful.")
}

func FetchAllMentors(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	var mentors []Mentor

	tx := db.
		Table("mentors").
		Select("name", "username").
		Find(&mentors)

	if tx.Error != nil {
		utils.LogErrAndRespond(r, w, tx.Error, "Database Error fetching mentors", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJson(r, w, mentors)
}

// /mentor/dashboard/ functions

func CreateMentorDashboard(mentor models.Mentor, db *gorm.DB) MentorDashboard {
	var projects []models.Project
	var projectsInfo []ProjectInfo
	var students []StudentInfo

	db.Table("projects").
		Where("mentor_id = ? OR secondary_mentor_id = ?", mentor.ID, mentor.ID).
		Select("name", "repo_link", "commit_count", "pull_count", "lines_added", "lines_removed", "contributors", "pulls", "project_status").
		Find(&projects)

	studentMap := make(map[string]bool)
	var studentUsernames []string
	for _, project := range projects {
		pulls := make([]string, 0)
		if len(project.Pulls) != 0 {
			pulls = strings.Split(project.Pulls, ",")
		}

		projectInfo := ProjectInfo{
			Name:          project.Name,
			RepoLink:      project.RepoLink,
			ProjectStatus: project.ProjectStatus,

			CommitCount:  project.CommitCount,
			PullCount:    project.PullCount,
			LinesAdded:   project.LinesAdded,
			LinesRemoved: project.LinesRemoved,

			Pulls: pulls,
		}
		projectsInfo = append(projectsInfo, projectInfo)

		for _, studentUsername := range strings.Split(project.Contributors, ",") {
			contains := studentMap[studentUsername]
			if contains {
				continue
			}

			studentMap[studentUsername] = true
			studentUsernames = append(studentUsernames, studentUsername)
		}
	}
	db.Table("students").Where("username IN ?", studentUsernames).
		Select("name", "username").Find(&students)

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
		Select("name", "username", "email", "ID").
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
