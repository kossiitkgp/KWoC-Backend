package controllers

import (
	"net/http"

	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
	"github.com/kossiitkgp/kwoc-backend/v2/models"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"
)

// FetchUnapprovedProjects godoc
//
//	@Summary		Fetches all unapproved Projects
//	@Description	Fetches the list of unapproved Projects
//	@Accept			plain
//	@Produce		json
//	@Success		200	{object}	[]Project	"Projects fetched successfully."
//	@Failure		500	{object}	utils.HTTPMessage	"Error fetching projects from the database."
//	@Router			/project/unapproved [get]
func FetchUnapprovedProjects(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	isAdmin := r.Context().Value(middleware.LOGIN_CTX_IS_ADMIN_KEY).(bool)
	if !isAdmin {
		utils.LogWarnAndRespond(r, w, "Error: Unauthorized access. Admins only.", http.StatusUnauthorized)
		return
	}

	var projects []models.Project

	tx := db.
		Table("projects").
		Preload("Mentor").
		Preload("SecondaryMentor").
		Where("project_status = ?", false).
		Select("id", "name", "description", "tags", "repo_link", "comm_channel", "readme_link", "mentor_id", "secondary_mentor_id").
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
