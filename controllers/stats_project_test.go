package controllers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kossiitkgp/kwoc-backend/controllers"
	"github.com/kossiitkgp/kwoc-backend/models"
)

func createFetchAllProjectStatsRequest() *http.Request {
	req, _ := http.NewRequest(
		"GET",
		"/stats/projects/",
		nil,
	)

	return req
}

func areProjectStatsEquivalent(stats *controllers.ProjectStats, project *models.Project) bool {
	return project.Name == stats.Name &&
		project.RepoLink == stats.RepoLink &&
		project.PullCount == stats.PullCount &&
		project.CommitCount == stats.CommitCount &&
		project.LinesAdded == stats.LinesAdded &&
		project.LinesRemoved == stats.LinesRemoved
}

func findProjIndex(repo_link string, list []models.Project) int {
	for i, proj := range list {
		if proj.RepoLink == repo_link {
			return i
		}
	}

	return 0
}

func TestFetchAllProjectStats(t *testing.T) {
	db := setTestDB()
	defer unsetTestDB()

	test_projects := generateTestProjects(10, true, true)

	req := createFetchAllProjectStatsRequest()
	res := executeRequest(req, db)

	var resStats []controllers.ProjectStats
	_ = json.NewDecoder(res.Body).Decode(&resStats)

	// Check if any projects with status false (not approved) are returned in the response
	var areAllProjectsApproved bool = true

	for _, stats := range resStats {
		testProj := test_projects[findProjIndex(stats.RepoLink, test_projects)]

		if !testProj.ProjectStatus {
			areAllProjectsApproved = false
			break
		}
	}

	if !areAllProjectsApproved {
		t.Fatalf("Unapproved projects (project_status = false) are returned by the /stats/projects/ endpoint.")
	}

	// Check if all the returned project stats and projects in the database are equal
	var areAllProjectsEquivalent bool = true

	for _, stats := range resStats {
		testProj := test_projects[findProjIndex(stats.RepoLink, test_projects)]

		if !areProjectStatsEquivalent(&stats, &testProj) {
			areAllProjectsEquivalent = false
			break
		}
	}

	if !areAllProjectsEquivalent {
		t.Fatalf("Stats returned by the /stats/projects/ endpoint are incorrect.")
	}
}
