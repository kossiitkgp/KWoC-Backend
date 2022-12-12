package controllers

import (
	"kwoc20-backend/utils"
	"net/http"
)

type OverallStatsRes struct {
	TotalCommits      uint `json:"total_commits"`
	TotalPRCount      uint `json:"total_pr_count"`
	TotalLinesAdded   uint `json:"total_lines_added"`
	TotalLinesRemoved uint `json:"total_lines_removed"`

	GenTime int64 `json:"generation_time"`
}

func OverallStats(req map[string]interface{}, r *http.Request) (interface{}, int) {
	db := utils.GetDB()

	var latest_stats OverallStatsRes

	result := db.
		Table("stats").
		Order("id").
		Last(&latest_stats)

	if result.Error != nil {
		return "Error: Database error", 500
	}

	return latest_stats, 200
}
