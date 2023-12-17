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

type UpdateMentorReqFields struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ProjectInfo struct {
	Id            uint     `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	RepoLink      string   `json:"repo_link"`
	ReadmeLink    string   `json:"readme_link"`
	Tags          []string `json:"tags"`
	ProjectStatus bool     `json:"project_status"`
	StatusRemark  string   `json:"status_remark"`

	CommitCount  uint `json:"commit_count"`
	PullCount    uint `json:"pull_count"`
	LinesAdded   uint `json:"lines_added"`
	LinesRemoved uint `json:"lines_removed"`

	Pulls           []string `json:"pulls"`
	Mentor          Mentor   `json:"mentor"`
	SecondaryMentor Mentor   `json:"secondary_mentor"`
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

	err = utils.DetectSessionHijackAndRespond(r, w, reqFields.Username, login_username, "Login username and given username do not match.")
	if err != nil {
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

	// Check if a student of the same username exists
	student := models.Student{}
	tx = db.
		Table("students").
		Where("username = ?", reqFields.Username).
		First(&student)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		utils.LogErrAndRespond(r, w, tx.Error, "Database error.", http.StatusInternalServerError)
		return
	}
	student_exists := student.Username == reqFields.Username

	if student_exists {
		utils.LogWarnAndRespond(
			r,
			w,
			fmt.Sprintf("The username `%s` already exists as a student.", reqFields.Username),
			http.StatusBadRequest,
		)

		return
	}

	// Create a db entry if the mentor doesn'tf exist
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
	var projects []models.Project = make([]models.Project, 0)
	var projectsInfo []ProjectInfo = make([]ProjectInfo, 0)
	var students []StudentInfo = make([]StudentInfo, 0)

	db.Preload("Mentor").Preload("SecondaryMentor").Table("projects").
		Where("mentor_id = ? OR secondary_mentor_id = ?", mentor.ID, mentor.ID).
		Find(&projects)

	studentMap := make(map[string]bool)
	var studentUsernames []string
	for _, project := range projects {
		pulls := make([]string, 0)
		if len(project.Pulls) != 0 {
			pulls = strings.Split(project.Pulls, ",")
		}

		tags := make([]string, 0)
		if len(project.Tags) != 0 {
			tags = strings.Split(project.Tags, ",")
		}

		projectInfo := ProjectInfo{
			Id:            project.ID,
			Name:          project.Name,
			Description:   project.Description,
			RepoLink:      project.RepoLink,
			ReadmeLink:    project.ReadmeLink,
			Tags:          tags,
			ProjectStatus: project.ProjectStatus,
			StatusRemark:  project.StatusRemark,

			CommitCount:  project.CommitCount,
			PullCount:    project.PullCount,
			LinesAdded:   project.LinesAdded,
			LinesRemoved: project.LinesRemoved,

			Pulls:           pulls,
			Mentor:          newMentor(&project.Mentor),
			SecondaryMentor: newMentor(&project.SecondaryMentor),
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

func UpdateMentorDetails(w http.ResponseWriter, r *http.Request) {
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

	var reqFields = UpdateMentorReqFields{}

	err := utils.DecodeJSONBody(r, &reqFields)
	if err != nil {
		utils.LogErrAndRespond(r, w, err, "Error decoding JSON body.", http.StatusBadRequest)
		return
	}

	tx = db.Model(&modelMentor).Updates(models.Mentor{
		Name:  reqFields.Name,
		Email: reqFields.Email,
	})

	if tx.Error != nil {
		utils.LogErrAndRespond(
			r,
			w,
			tx.Error,
			"Invalid Details: Could not update mentor details",
			http.StatusBadRequest,
		)
		return
	}

	utils.RespondWithJson(r, w, []string{"Mentor details updated successfully."})
}

func GetMentorDetails(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	login_username := r.Context().Value(middleware.LoginCtxKey(middleware.LOGIN_CTX_USERNAME_KEY))

	mentor := models.Mentor{}
	tx := db.
		Table("mentors").
		Where("username = ?", login_username).
		Select("name", "username", "email", "ID").
		First(&mentor)

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

	utils.RespondWithJson(r, w, mentor)
}
