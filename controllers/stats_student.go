package controllers

import (
	"kwoc-backend/middleware"
	"kwoc-backend/models"
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

	var students []models.Student

	tx := db.
		Table("students").
		Select("name", "username", "pull_count", "commit_count", "lines_added", "lines_removed").
		Find(&students)

	if tx.Error != nil {
		utils.LogErrAndRespond(
			r,
			w,
			tx.Error,
			"Error fetching projects from the database.",
			http.StatusInternalServerError,
		)
		return
	}

	var response []StudentBriefStats = make([]StudentBriefStats, 0)

	for _, student := range students {
		response = append(
			response,
			StudentBriefStats{
				Name:         student.Name,
				Username:     student.Username,
				PullCount:    student.PullCount,
				CommitCount:  student.CommitCount,
				LinesAdded:   student.LinesAdded,
				LinesRemoved: student.LinesRemoved,
			},
		)
	}

	utils.RespondWithJson(r, w, response)
}
