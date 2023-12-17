package controllers_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kossiitkgp/kwoc-backend/v2/controllers"
	"github.com/kossiitkgp/kwoc-backend/v2/models"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"
)

func createFetchAllProjRequest() *http.Request {
	req, _ := http.NewRequest(
		"GET",
		"/project/",
		nil,
	)

	return req
}

func createFetchProjDetailsRequest(id any) *http.Request {
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("/project/%v", id),
		nil,
	)

	return req
}

func generateTestProjects(numProjects int, randomizeProjectStatus bool, defaultProjectStatus bool) []models.Project {
	rand.Seed(time.Now().Unix())

	var projects []models.Project = make([]models.Project, 0)

	for i := 0; i < numProjects; i++ {
		var projectStatus bool

		if randomizeProjectStatus {
			projectStatus = rand.Intn(10) > 5
		} else {
			projectStatus = defaultProjectStatus
		}

		projects = append(
			projects,
			models.Project{
				Name:          fmt.Sprintf("YANGJF-%d", rand.Int()),
				Description:   fmt.Sprintf("Yet another next-gen javascript framework v%d.1", rand.Int()),
				Tags:          fmt.Sprintf("next-gen, javascript, framework, %dth iteration", rand.Int()),
				RepoLink:      "https://xkcd.com/927/",
				CommChannel:   fmt.Sprintf("https://link%d", rand.Int()),
				ReadmeLink:    fmt.Sprintf("https://readme%d", rand.Int()),
				ProjectStatus: projectStatus,
				StatusRemark:  fmt.Sprintf("Status remark %d", rand.Int()),

				// Stats
				CommitCount:  uint(rand.Int()),
				PullCount:    uint(rand.Int()),
				LinesAdded:   uint(rand.Int()),
				LinesRemoved: uint(rand.Int()),
			},
		)
	}

	return projects
}

func areProjectsEquivalent(proj1 *controllers.Project, proj2 *models.Project) bool {
	return proj1.Name == proj2.Name &&
		proj1.Description == proj2.Description &&
		strings.Join(proj1.Tags, ",") == proj2.Tags &&
		proj1.RepoLink == proj2.RepoLink &&
		proj1.CommChannel == proj2.CommChannel &&
		proj1.ReadmeLink == proj2.ReadmeLink
}

func TestFetchAllProjects(t *testing.T) {
	db := setTestDB()
	defer unsetTestDB()

	testProjects := generateTestProjects(10, true, true)

	_ = db.Table("projects").Create(testProjects)

	req := createFetchAllProjRequest()
	res := executeRequest(req, db)

	var resProjects []controllers.Project
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
			break
		}
	}

	if !areAllProjectsEquivalent {
		t.Fatalf("Projects returned by the /project/ endpoint are incorrect.")
	}
}

// Try fetching a project with an invalid id
func TestFetchProjDetailsInvalidID(t *testing.T) {
	req := createFetchProjDetailsRequest("kekw")
	res := executeRequest(req, nil)

	expectStatusCodeToBe(t, res, http.StatusBadRequest)
	expectResponseJSONBodyToBe(t, res, utils.HTTPMessage{StatusCode: http.StatusBadRequest, Message: "Error parsing project id."})
}

// Try fetching a project that does not exist
func TestFetchProjDetailsDNE(t *testing.T) {
	db := setTestDB()
	defer unsetTestDB()

	testProjId := rand.Int()

	req := createFetchProjDetailsRequest(testProjId)
	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusBadRequest)
	expectResponseJSONBodyToBe(t, res, utils.HTTPMessage{StatusCode: http.StatusBadRequest, Message: fmt.Sprintf("Project with id `%d` does not exist.", testProjId)})
}

// Try to fetch an unapproved project
func TestFetchProjDetailsUnapproved(t *testing.T) {
	db := setTestDB()
	defer unsetTestDB()

	testProj := generateTestProjects(1, false, false)[0]

	_ = db.Table("projects").Create(&testProj)

	req := createFetchProjDetailsRequest(1)
	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusBadRequest)
	expectResponseJSONBodyToBe(t, res, utils.HTTPMessage{StatusCode: http.StatusBadRequest, Message: fmt.Sprintf("Project with id `%d` does not exist.", 1)})
}

// Try to fetch a valid project
func TestFetchProjDetailsOK(t *testing.T) {
	db := setTestDB()
	defer unsetTestDB()

	testProjects := generateTestProjects(5, false, true)

	_ = db.Table("projects").Create(testProjects)

	for i, proj := range testProjects {
		req := createFetchProjDetailsRequest(i + 1)
		res := executeRequest(req, db)

		var resProj controllers.Project

		_ = json.NewDecoder(res.Body).Decode(&resProj)

		if !areProjectsEquivalent(&resProj, &proj) {
			t.Fatalf("/projects/%d returned incorrect information.", i+1)
		}
	}
}
