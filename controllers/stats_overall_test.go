package controllers_test

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/kossiitkgp/kwoc-backend/v2/controllers"
	"github.com/kossiitkgp/kwoc-backend/v2/models"
)

func createFetchOverallStatsRequest() *http.Request {
	req, _ := http.NewRequest(
		"GET",
		"/stats/overall/",
		nil,
	)

	return req
}

func generateTestStats(numStatsEntries int) []models.Stats {
	rand.Seed(time.Now().Unix())

	var stats []models.Stats = make([]models.Stats, 0)

	stats = append(
		stats,
		models.Stats{
			GenTime:           time.Now().Unix(),
			TotalCommitCount:  uint(rand.Int()),
			TotalPullCount:    uint(rand.Int()),
			TotalLinesAdded:   uint(rand.Int()),
			TotalLinesRemoved: uint(rand.Int()),
		},
	)

	for i := 1; i < numStatsEntries; i++ {
		prev_entry := stats[i-1]
		stats = append(
			stats,
			models.Stats{
				GenTime:           prev_entry.GenTime + int64(rand.Intn(10)),
				TotalCommitCount:  prev_entry.TotalCommitCount + uint(rand.Int()),
				TotalPullCount:    prev_entry.TotalPullCount + uint(rand.Int()),
				TotalLinesAdded:   prev_entry.TotalLinesAdded + uint(rand.Int()),
				TotalLinesRemoved: prev_entry.TotalLinesRemoved + uint(rand.Int()),
			},
		)
	}

	return stats
}

func areOverallStatsEquivalent(stats *controllers.OverallStats, dbStats *models.Stats) bool {
	return stats.GenTime == dbStats.GenTime &&
		stats.TotalCommitCount == dbStats.TotalCommitCount &&
		stats.TotalLinesAdded == dbStats.TotalLinesAdded &&
		stats.TotalLinesRemoved == dbStats.TotalLinesRemoved &&
		stats.TotalPullCount == dbStats.TotalPullCount
}

func TestFetchOverallStats(t *testing.T) {
	db := setTestDB()
	defer unsetTestDB()

	testStats := generateTestStats(10)

	_ = db.Table("stats").Create(testStats)

	req := createFetchOverallStatsRequest()
	res := executeRequest(req, db)

	var resStats []controllers.OverallStats
	_ = json.NewDecoder(res.Body).Decode(&resStats)

	var areAllStatsEquivalent bool = true

	for i, entry := range resStats {
		// Assuming stats are created in order
		testEntry := testStats[i]

		if !areOverallStatsEquivalent(&entry, &testEntry) {
			areAllStatsEquivalent = false
			break
		}
	}

	if !areAllStatsEquivalent {
		t.Fatalf("Stats returned by the /stats/overall/ endpoint are incorrect.")
	}
}
