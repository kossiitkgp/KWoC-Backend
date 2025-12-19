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

/*
	Public-safe mentor response
	Only username + display name
	(Required for Issue #177)
*/
type PublicMentor struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
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

	Pulls           []string      `json:"pulls"`
	Mentor          PublicMentor  `json:"mentor"`
	SecondaryMentor PublicMentor  `json:"secondary_mentor"`
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

/*
	Helper: Convert Mentor → PublicMentor
*/
func newMentor(m *models.Mentor) PublicMentor {
	if m == nil {
		return PublicMentor{}
	}

	return PublicMentor{
		Username:    m.Username,
		DisplayName: m.Name,
	}
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

	loginUsername := r.Context().Value(middleware.LOGIN_CTX_USERNAME_KEY).(string)

	err = utils.DetectSessionHijackAndRespond(
		r,
		w,
		reqFields.Username,
		loginUsername,
		"Login username and given username do not match.",
	)
	if err != nil {
		return
	}

	mentor := models.Mentor{}
	tx := db.Table("mentors").
		Where("username = ?", reqFields.Username).
		First(&mentor)

	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		utils.LogErrAndRespond(r, w, tx.Error, "Database error.", http.StatusInternalServerError)
		return
	}

	if mentor.Username == reqFields.Username {
		utils.LogWarnAndRespond(
			r,
			w,
			fmt.Sprintf("Mentor `%s` already exists.", mentor.Username),
			http.StatusBadRequest,
		)
		return
	}

	student := models.Student{}
	tx = db.Table("students").
		Where("username = ?", reqFields.Username).
		First(&student)

	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		utils.LogErrAndRespond(r, w, tx.Error, "Database error.", http.StatusInternalServerError)
		return
	}

	if student.Username == reqFields.Username {
		utils.LogWarnAndRespond(
			r,
			w,
			fmt.Sprintf("The username `%s` already exists as a student.", reqFields.Username),
			http.StatusBadRequest,
		)
		return
	}

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

func CreateMentorDashboard(mentor models.Mentor, db *gorm.DB) MentorDashboard {
	var projects []models.Project
	var projectsInfo []ProjectInfo
	var students []StudentInfo

	db.Preload("Mentor").
		Preload("SecondaryMentor").
		Table("projects").
		Where("mentor_id = ? OR secondary_mentor_id = ?", mentor.ID, mentor.ID).
		Find(&projects)

	studentMap := make(map[string]bool)
	var studentUsernames []string

	for _, project := range projects {
		pulls := []string{}
		if project.Pulls != "" {
			pulls = strings.Split(project.Pulls, ",")
		}

		tags := []string{}
		if project.Tags != "" {
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
			if studentUsername == "" || studentMap[studentUsername] {
				continue
			}
			studentMap[studentUsername] = true
			studentUsernames = append(studentUsernames, studentUsername)
		}
	}

	db.Table("students").
		Where("username IN ?", studentUsernames).
		Select("name", "username").
		Find(&students)

	return MentorDashboard{
		Name:     mentor.Name,
		Username: mentor.Username,
		Email:    mentor.Email,
		Projects: projectsInfo,
		Students: students,
	}
}
