package controllers_test

import (
	"encoding/json"
	"fmt"
	"kwoc-backend/controllers"
	"kwoc-backend/models"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func createFetchAllProjRequest() *http.Request {
	req, _ := http.NewRequest(
		"GET",
		"/projects/",
		nil,
	)

	return req
}

func createFetchProjDetailsRequest(id uint) *http.Request {
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("/project/%d", id),
		nil,
	)

	return req
}

func generateTestProjects(numProjects int, randomizeProjectStatus bool) []models.Project {
	rand.Seed(time.Now().Unix())

	var projects []models.Project = make([]models.Project, 0)

	for i := 0; i < numProjects; i++ {
		var projectStatus bool

		if randomizeProjectStatus {
			projectStatus = rand.Intn(10) > 5
		} else {
			projectStatus = true
		}

		projects = append(
			projects,
			models.Project{
				Name:          fmt.Sprintf("YANGJF-%d", rand.Int()),
				Desc:          fmt.Sprintf("Yet another next-gen javascript framework v%d.1", rand.Int()),
				Tags:          fmt.Sprintf("next-gen, javascript, framework, %dth iteration", rand.Int()),
				RepoLink:      "https://xkcd.com/927/",
				ComChannel:    fmt.Sprintf("https://link%d", rand.Int()),
				README:        fmt.Sprintf("https://readme%d", rand.Int()),
				ProjectStatus: projectStatus,
			},
		)
	}

	return projects
}

func areProjectsEquivalent(proj1 *controllers.FetchProjProject, proj2 *models.Project) bool {
	return proj1.Name == proj2.Name &&
		proj1.Desc == proj2.Desc &&
		proj1.Tags == proj2.Tags &&
		proj1.RepoLink == proj2.RepoLink &&
		proj1.ComChannel == proj2.ComChannel &&
		proj1.ReadmeURL == proj2.README
}

func TestFetchAllProjects(t *testing.T) {
	db := setTestDB()
	defer unsetTestDB()

	testProjects := generateTestProjects(10, true)

	db.Table("projects").Create(testProjects)

	req := createFetchAllProjRequest()
	res := executeRequest(req, db)

	var resProjects []controllers.FetchProjProject
	_ = json.NewDecoder(res.Body).Decode(&resProjects)

	// Check if any projects with status false (not approved) are returned in the request
	var areAllProjectsApproved bool = true

	for _, proj := range resProjects {
		// Assuming projects are created in order
		testProj := testProjects[proj.Id-1]

		if !testProj.ProjectStatus {
			areAllProjectsApproved = false
			break
		}
	}

	if !areAllProjectsApproved {
		t.Fatalf("Unapproved projects (project_status = false) are returned by the /project/ endpoint.")
	}

	// Check if all the returned projects and projects in the database are equal
	var areAllProjectsEquivalent bool = true

	for _, proj := range resProjects {
		// Assuming projects are created in order
		testProj := testProjects[proj.Id-1]

		if !areProjectsEquivalent(&proj, &testProj) {
			areAllProjectsEquivalent = false
			return
		}
	}

	if !areAllProjectsEquivalent {
		t.Fatalf("Projects returned by the /project/ endpoint are incorrect.")
	}
}
