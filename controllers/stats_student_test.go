package controllers_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/kossiitkgp/kwoc-backend/v2/controllers"
	"github.com/kossiitkgp/kwoc-backend/v2/models"
)

func createFetchAllStudentsStatsRequest() *http.Request {
	req, _ := http.NewRequest(
		"GET",
		"/stats/students/",
		nil,
	)

	return req
}

func generateTestStudents(numStudents int) []models.Student {
	rand.Seed(time.Now().Unix())

	var students []models.Student = make([]models.Student, 0)

	for i := 0; i < numStudents; i++ {
		students = append(
			students,
			models.Student{
				Name:         fmt.Sprintf("TestStudent-%d", rand.Int()),
				Username:     fmt.Sprintf("TestUsername-%d", rand.Int()),
				PullCount:    uint(rand.Uint32()),
				CommitCount:  uint(rand.Uint32()),
				LinesAdded:   uint(rand.Uint32()),
				LinesRemoved: uint(rand.Uint32()),
			},
		)
	}

	return students
}

func areStudentStatsEquivalent(stats *controllers.StudentBriefStats, student *models.Student) bool {
	return student.Name == stats.Name &&
		student.Username == stats.Username &&
		student.PullCount == stats.PullCount &&
		student.CommitCount == stats.CommitCount &&
		student.LinesAdded == stats.LinesAdded &&
		student.LinesRemoved == stats.LinesRemoved
}

func TestFetchAllStudentsStats(t *testing.T) {
	db := setTestDB()
	defer unsetTestDB()

	testStudents := generateTestStudents(10)

	_ = db.Table("students").Create(testStudents)

	req := createFetchAllStudentsStatsRequest()
	res := executeRequest(req, db)

	var resStats []controllers.StudentBriefStats
	_ = json.NewDecoder(res.Body).Decode(&resStats)

	// Check if all the returned stats and students in the database are equal
	var areAllStudentsEquivalent bool = true

	for i, stats := range resStats {
		// Assuming students are created in order
		testStudent := testStudents[i]

		if !areStudentStatsEquivalent(&stats, &testStudent) {
			areAllStudentsEquivalent = false
			break
		}
	}

	if !areAllStudentsEquivalent {
		t.Fatalf("Stats returned by the /stats/students endpoint are incorrect.")
	}
}
