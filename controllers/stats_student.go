package controllers

import (
	"kwoc-backend/middleware"
	"kwoc-backend/utils"
	"net/http"
)

type StudentBriefStats struct {
	Name     string `json:"name"`
	Username string `json:"username"`

	PullCount    uint `json:"pull_count"`
	CommitCount  uint `json:"commit_count"`
	LinesAdded   uint `json:"lines_added"`
	LinesRemoved uint `json:"lines_removed"`
}

func FetchAllStudentStats(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	var student_stats []StudentBriefStats

	tx := db.
		Table("students").
		Select("name", "username", "pull_count", "commit_count", "lines_added", "lines_removed").
		Find(&student_stats)

	if tx.Error != nil {
		utils.LogErrAndRespond(
			r,
			w,
			tx.Error,
			"Error fetching student stats from the database.",
			http.StatusInternalServerError,
		)
		return
	}

	utils.RespondWithJson(r, w, student_stats)
}
