package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
	"github.com/kossiitkgp/kwoc-backend/v2/models"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"
)

// ApproveProject godoc
//
//	@Summary		Approves a project
//	@Description	Approves a project given its ID
//	@Accept			plain
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Project approved successfully."
//	@Failure		500	{object}	utils.HTTPMessage	"Error approving project."
//	@Router			/project/{id}/approve [post]
func ApproveProject(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	isAdmin := r.Context().Value(middleware.LOGIN_CTX_IS_ADMIN_KEY).(bool)
	if !isAdmin {
		utils.LogWarnAndRespond(r, w, "Error: Unauthorized access. Admins only.", http.StatusUnauthorized)
		return
	}

	reqParams := mux.Vars(r)
	if reqParams["id"] == "" {
		utils.LogWarnAndRespond(r, w, "Error: Project ID not provided in the request URL.", http.StatusBadRequest)
		return
	}
	projectIdStr := reqParams["id"]

	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		utils.LogErrAndRespond(r, w, err, "Error parsing project ID.", http.StatusBadRequest)
		return
	}

	project := models.Project{}
	tx := db.
		Table("projects").
		Where("id = ?", projectId).
		First(&project)

	if tx.Error != nil {
		utils.LogErrAndRespond(r, w, tx.Error, "Error fetching project from the database.", http.StatusInternalServerError)
		return
	}

	if project.ProjectStatus {
		utils.LogWarnAndRespond(r, w, "Error: Project is already approved.", http.StatusBadRequest)
		return
	}

	project.ProjectStatus = true
	saveTx := db.Save(&project)
	if saveTx.Error != nil {
		utils.LogErrAndRespond(r, w, saveTx.Error, "Error updating project status in the database.", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJson(r, w, map[string]string{"message": "Project approved successfully."})
}
