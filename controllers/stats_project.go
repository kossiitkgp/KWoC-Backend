package controllers

import (
	"net/http"

	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"
)

type ProjectStats struct {
	Name     string `json:"name"`
	RepoLink string `json:"repo_link"`

	CommitCount  uint `json:"commit_count"`
	PullCount    uint `json:"pull_count"`
	LinesAdded   uint `json:"lines_added"`
	LinesRemoved uint `json:"lines_removed"`
}

func FetchAllProjectStats(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	var project_stats []ProjectStats

	tx := db.
		Table("projects").
		Where("project_status = ?", true).
		Select("name", "repo_link", "commit_count", "pull_count", "lines_added", "lines_removed").
		Find(&project_stats)

	if tx.Error != nil {
		utils.LogErrAndRespond(
			r,
			w,
			tx.Error,
			"Error fetching project stats from the database.",
			http.StatusInternalServerError,
		)
		return
	}

	utils.RespondWithJson(r, w, project_stats)
}
