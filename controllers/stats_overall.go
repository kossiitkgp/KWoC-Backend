package controllers

import (
	"net/http"

	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"
)

type OverallStats struct {
	TotalCommitCount  uint `json:"total_commit_count"`
	TotalPullCount    uint `json:"total_pull_count"`
	TotalLinesAdded   uint `json:"total_lines_added"`
	TotalLinesRemoved uint `json:"total_lines_removed"`

	GenTime int64 `json:"gen_time"`
}

func FetchOverallStats(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	var overall_stats []OverallStats

	tx := db.
		Table("stats").
		Order("gen_time DESC").
		Find(&overall_stats)

	if tx.Error != nil {
		utils.LogErrAndRespond(
			r,
			w,
			tx.Error,
			"Error fetching stats from the database.",
			http.StatusInternalServerError,
		)
		return
	}

	utils.RespondWithJson(r, w, overall_stats)
}
