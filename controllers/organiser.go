package controllers

import (
	"fmt"
	"net/http"

	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
	"github.com/kossiitkgp/kwoc-backend/v2/models"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"
)

type AcceptRejectProject struct {
	// Id of the project in the database (required)
	Id uint `json:"id"`
	// Status to be set of the project
	ProjectStatus bool `json:"project_status"`
	// Status Remark to be set of the project
	StatusRemark string `json:"status_remark"`
}

func OrgFetchAllProjectDetails(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db
	user_details := r.Context().Value(middleware.LOGIN_CTX_USERNAME_KEY).(utils.LoginJwtFields)

	if user_details.UserType != "organiser" {
		utils.LogErrAndRespond(r, w, nil, fmt.Sprintf("Error '%s' is not an organiser", user_details.Username), 400)
		return
	}

	var projects []models.Project

	tx := db.
		Table("projects").
		Preload("Mentor").
		Preload("SecondaryMentor").
		Select("id", "name", "description", "tags", "repo_link", "comm_channel", "readme_link", "mentor_id", "secondary_mentor_id", "project_status", "status_remark", "pull_count").
		Find(&projects)

	if tx.Error != nil {
		utils.LogErrAndRespond(r, w, tx.Error, "Error fetching projects from the database.", http.StatusInternalServerError)
		return
	}

	var response []Project = make([]Project, 0)

	for _, project := range projects {
		response = append(response, newProject(&project))
	}

	utils.RespondWithJson(r, w, response)
}

func UpdateStatusProject(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db
	user_details := r.Context().Value(middleware.LOGIN_CTX_USERNAME_KEY).(utils.LoginJwtFields)

	if user_details.UserType != OAUTH_TYPE_ORGANISER {
		utils.LogErrAndRespond(r, w, nil, fmt.Sprintf("Error '%s' is not an organiser", user_details.Username), 400)
		return
	}

	projectDetails := AcceptRejectProject{}

	err := utils.DecodeJSONBody(r, &projectDetails)
	if err != nil {
		utils.LogErrAndRespond(r, w, err, "Error decoding request JSON body.", http.StatusBadRequest)
		return
	}

	tx := db.
		Table("projects").
		Where("id = ?", projectDetails.Id).
		Updates(projectDetails)

	if tx.Error != nil {
		utils.LogErrAndRespond(r, w, tx.Error, "Error updating the project.", http.StatusInternalServerError)
		return
	}

	if projectDetails.ProjectStatus {
		utils.RespondWithHTTPMessage(r, w, http.StatusOK, "Project accepted successfully.")
	} else {
		utils.RespondWithHTTPMessage(r, w, http.StatusOK, "Project Rejected successfully.")
	}

}
