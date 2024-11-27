package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
	"github.com/kossiitkgp/kwoc-backend/v2/models"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"

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

// UpdateProject godoc
//
// @Summary		Update Project Details
// @Description	Update project details for the provided project ID.
// @Accept			json
// @Produce		json
// @Param			request	body		UpdateProjectReqFields	true	"Fields required for Project update."
// @Success		200		{object}	utils.HTTPMessage	"Project successfully updated."
// @Failure		401		{object}	utils.HTTPMessage	"Login username and mentor username do not match."
// @Failure		400		{object}	utils.HTTPMessage	"Error decoding request JSON body."
// @Failure		400		{object}	utils.HTTPMessage	"Mentor `username` does not exists."
// @Failure		400		{object}	utils.HTTPMessage	"Invalid Details: Could not update mentor details"
// @Failure		400		{object}	utils.HTTPMessage	"Error: Project `repo_link` does not exist."
// @Failure		400		{object}	utils.HTTPMessage	"Error: Mentor `username` does not own the project with ID `id`."
// @Failure		400		{object}	utils.HTTPMessage	"Error: Secondary mentor `secondary_mentor_username` cannot be same as primary mentor."
// @Failure		500		{object}	utils.HTTPMessage	"Error updating the project."
// @Failure		500		{object}	utils.HTTPMessage	"Database error."
// @Failure		500		{object}	utils.HTTPMessage	"Error fetching secondary mentor `secondary_mentor_username`."
//
// @Security		JWT
//
// @Router			/project/ [put]
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	reqFields := UpdateProjectReqFields{}

	err := utils.DecodeJSONBody(r, &reqFields)
	if err != nil {
		utils.LogErrAndRespond(r, w, err, "Error decoding request JSON body.", http.StatusBadRequest)
		return
	}

	login_username := r.Context().Value(middleware.LoginCtxKey(middleware.LOGIN_CTX_USERNAME_KEY)).(utils.LoginJwtFields).Username

	err = utils.DetectSessionHijackAndRespond(r, w, reqFields.MentorUsername, login_username, "Login username and mentor username do not match.")
	if err != nil {
		return
	}

	// Check if the project already exists
	project := models.Project{}
	tx := db.
		Table("projects").
		Preload("Mentor").
		Preload("SecondaryMentor").
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
		if reqFields.MentorUsername == reqFields.SecondaryMentorUsername {
			utils.LogErrAndRespond(
				r,
				w,
				err,
				fmt.Sprintf("Error: Secondary mentor `%s` cannot be same as primary mentor.", reqFields.SecondaryMentorUsername),
				http.StatusBadRequest,
			)
			return
		}

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
		}
	}

	updatedProj := &models.Project{
		Name:            reqFields.Name,
		Description:     reqFields.Description,
		Tags:            strings.Join(reqFields.Tags, ","),
		RepoLink:        reqFields.RepoLink,
		CommChannel:     reqFields.CommChannel,
		ReadmeLink:      reqFields.ReadmeLink,
		SecondaryMentor: secondaryMentor,
	}

	tx = db.
		Table("projects").
		Where("id = ?", reqFields.Id).
		Updates(updatedProj)

	if tx.Error != nil {
		utils.LogErrAndRespond(r, w, tx.Error, "Error updating the project.", http.StatusInternalServerError)
		return
	}

	utils.RespondWithHTTPMessage(r, w, http.StatusOK, "Project successfully updated.")
}
