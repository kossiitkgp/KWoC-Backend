package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
	"github.com/kossiitkgp/kwoc-backend/v2/models"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UpdateProjectReqFields struct {
	// Id of the project in the database (required)
	Id uint `json:"id"`
	// Name of the project
	Name string `json:"name"`
	// Description for the project
	Description string `json:"description"`
	// List of tags for the project
	Tags []string `json:"tags"`
	// Mentor's username
	MentorUsername string `json:"mentor_username"`
	// Secondary mentor's username (if updated)
	SecondaryMentorUsername string `json:"secondary_mentor_username"`
	// Link to the repository of the project
	RepoLink string `json:"repo_link"`
	// Link to a communication channel/platform
	CommChannel string `json:"comm_channel"`
	// Link to the project's README file
	ReadmeLink string `json:"readme_link"`
}

func UpdateProject(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	reqFields := UpdateProjectReqFields{}

	err := utils.DecodeJSONBody(r, &reqFields)
	if err != nil {
		utils.LogErrAndRespond(r, w, err, "Error decoding request JSON body.", http.StatusBadRequest)
		return
	}

	login_username := r.Context().Value(middleware.LoginCtxKey(middleware.LOGIN_CTX_USERNAME_KEY))

	if reqFields.MentorUsername != login_username {
		log.Warn().Msgf(
			"%s %s %s\n%s %s",
			r.Method,
			r.RequestURI,
			"POSSIBLE SESSION HIJACKING.",
			fmt.Sprintf("JWT Username: %s", login_username),
			fmt.Sprintf("Given Username: %s", reqFields.MentorUsername),
		)

		utils.RespondWithHTTPMessage(r, w, http.StatusUnauthorized, "Login username and mentor username do not match.")
		return
	}

	// Check if the project already exists
	project := models.Project{}
	tx := db.
		Table("projects").
		Where("id = ?", reqFields.Id).
		First(&project)

	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		utils.LogErrAndRespond(r, w, err, "Database error.", http.StatusInternalServerError)
		return
	}

	if tx.Error == gorm.ErrRecordNotFound {
		utils.LogErrAndRespond(
			r,
			w,
			tx.Error,
			fmt.Sprintf("Error: Project `%s` does not exist.", reqFields.RepoLink),
			http.StatusBadRequest,
		)
		return
	}

	if project.Mentor.Username != login_username {
		utils.LogErrAndRespond(
			r,
			w,
			tx.Error,
			fmt.Sprintf("Error: Mentor `%s` does not own the project with ID `%d`.", login_username, project.ID),
			http.StatusBadRequest,
		)
		return
	}

	// Attempt to fetch secondary mentor from the database
	secondaryMentor := models.Mentor{}
	if reqFields.SecondaryMentorUsername != "" {
		tx = db.Table("mentors").Where("username = ?", reqFields.SecondaryMentorUsername).First(&secondaryMentor)

		if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
			utils.LogErrAndRespond(
				r,
				w,
				err,
				fmt.Sprintf("Error fetching secondary mentor `%s`.", reqFields.SecondaryMentorUsername),
				http.StatusInternalServerError,
			)
			return
		} else if tx.Error == gorm.ErrRecordNotFound {
			utils.LogWarnAndRespond(
				r,
				w,
				fmt.Sprintf("Secondary mentor `%s` does not exist.", reqFields.SecondaryMentorUsername),
				http.StatusBadRequest,
			)
			return
		}
	}

	secondaryMentorId := int32(secondaryMentor.ID)
	updatedProj := &models.Project{
		Name:              reqFields.Name,
		Description:       reqFields.Description,
		Tags:              strings.Join(reqFields.Tags, ","),
		RepoLink:          reqFields.RepoLink,
		CommChannel:       reqFields.CommChannel,
		ReadmeLink:        reqFields.ReadmeLink,
		SecondaryMentorId: &secondaryMentorId,
	}

	tx = db.
		Table("projects").
		Where("id = ?", reqFields.Id).
		Select("name", "description", "tags", "repo_link", "comm_channel", "readme_link", "secondary_mentor_id").
		Updates(updatedProj)

	if tx.Error != nil {
		utils.LogErrAndRespond(r, w, tx.Error, "Error updating the project.", http.StatusInternalServerError)
		return
	}

	utils.RespondWithHTTPMessage(r, w, http.StatusOK, "Project successfully updated.")
}
