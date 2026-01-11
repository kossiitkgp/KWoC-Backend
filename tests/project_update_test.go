package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"testing"

	"github.com/kossiitkgp/kwoc-backend/v2/controllers"
	"github.com/kossiitkgp/kwoc-backend/v2/models"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"

	"gorm.io/gorm"
)

func createProjectUpdateRequest(reqFields *controllers.UpdateProjectReqFields) *http.Request {
	reqBody, _ := json.Marshal(reqFields)

	req, _ := http.NewRequest(
		"PUT",
		"/project/",
		bytes.NewReader(reqBody),
	)

	return req
}

// Test unauthenticated request to /project/
func TestProjectUpdateNoAuth(t *testing.T) {
	testRequestNoAuth(t, "PUT", "/project/")
}

// Test request to /project/ with invalid jwt
func TestProjectUpdateInvalidAuth(t *testing.T) {
	testRequestInvalidAuth(t, "PUT", "/project/")
}

// Test request to /project/ with session hijacking attempt
func TestProjectUpdateSessionHijacking(t *testing.T) {
	// Generate a jwt secret key for testing
	setTestJwtSecretKey()
	defer unsetTestJwtSecretKey()

	testLoginFields := utils.LoginJwtFields{Username: "someuser"}

	someuserJwt, _ := utils.GenerateLoginJwtString(testLoginFields)

	reqFields := controllers.UpdateProjectReqFields{Id: 1, MentorUsername: "anotheruser"}

	req := createProjectUpdateRequest(&reqFields)
	req.Header.Add("Bearer", someuserJwt)

	res := executeRequest(req, nil)

	expectStatusCodeToBe(t, res, http.StatusUnauthorized)
	expectResponseJSONBodyToBe(t, res, utils.HTTPMessage{StatusCode: http.StatusUnauthorized, Message: "Login username and mentor username do not match."})
}

// Test a request to /project/ to update a non-existent project
func tProjectUpdateNonExistent(db *gorm.DB, testUsername string, testJwt string, t *testing.T) {
	projectReqFields := &controllers.UpdateProjectReqFields{
		Id:             2,
		RepoLink:       "https://example.com",
		MentorUsername: testUsername,
	}

	projectReq := createProjectUpdateRequest(projectReqFields)
	projectReq.Header.Add("Bearer", testJwt)

	projectRes := executeRequest(projectReq, db)

	expectStatusCodeToBe(t, projectRes, http.StatusBadRequest)
	expectResponseJSONBodyToBe(t, projectRes, utils.HTTPMessage{StatusCode: http.StatusBadRequest, Message: fmt.Sprintf("Error: Project `%s` does not exist.", projectReqFields.RepoLink)})
}

// Test a request to /project/ to update an existent project
func tProjectUpdateExistent(db *gorm.DB, testUsername string, testJwt string, t *testing.T) {
	// Register a test project
	projRegFields := createTestProjectRegFields(testUsername, "")
	proj := models.Project{
		Name:          projRegFields.Name,
		Description:   projRegFields.Description,
		Tags:          strings.Join(projRegFields.Tags, ","),
		RepoLink:      projRegFields.RepoLink,
		CommChannel:   projRegFields.CommChannel,
		
		ProjectStatus: true,

		Mentor: models.Mentor{
			Username: projRegFields.MentorUsername,
		},
		SecondaryMentor: models.Mentor{
			Username: projRegFields.SecondaryMentorUsername,
		},
	}
	db.Table("projects").Create(&proj)

	// Create updated fields
	projUpdateFields := &controllers.UpdateProjectReqFields{
		Id:             proj.ID,
		Name:           fmt.Sprintf("Nename %d", rand.Int()),
		Description:    "New description.",
		Tags:           strings.Split("New,Tags,Test", ","),
		MentorUsername: testUsername,
		RepoLink:       "http://NewRepoLink",
		CommChannel:    "totallynewcomchannel",
		
	}

	// Test with a valid new secondary mentor
	projUpdateFields.SecondaryMentorUsername = "testSecondary"

	req := createProjectUpdateRequest(projUpdateFields)
	req.Header.Add("Bearer", testJwt)

	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusOK)
	expectResponseJSONBodyToBe(t, res, utils.HTTPMessage{StatusCode: http.StatusOK, Message: "Project successfully updated."})

	// Check if the project got updated
	var updatedProj models.Project
	_ = db.
		Table("projects").
		Preload("Mentor").
		Preload("SecondaryMentor").
		Where("id = ?", projUpdateFields.Id).
		First(&updatedProj)

	if updatedProj.Name != projUpdateFields.Name {
		t.Errorf("Project Name field did not get updated\n Expected: `%s`. Received: `%s`", projUpdateFields.Name, updatedProj.Name)
	}

	if updatedProj.Description != projUpdateFields.Description {
		t.Errorf("Project Description field did not get updated\n Expected: `%s`. Received: `%s`", projUpdateFields.Description, updatedProj.Description)
	}

	if updatedProj.Tags != strings.Join(projUpdateFields.Tags, ",") {
		t.Errorf("Project Tags field did not get updated\n Expected: `%s`. Received: `%s`", projUpdateFields.Tags, updatedProj.Tags)
	}

	if updatedProj.RepoLink != projUpdateFields.RepoLink {
		t.Errorf("Project RepoLink field did not get updated\n Expected: `%s`. Received: `%s`", projUpdateFields.RepoLink, updatedProj.RepoLink)
	}

	if updatedProj.CommChannel != projUpdateFields.CommChannel {
		t.Errorf("Project CommChannel field did not get updated\n Expected: `%s`. Received: `%s`", projUpdateFields.CommChannel, updatedProj.CommChannel)
	}

	

	if updatedProj.SecondaryMentor.Username != projUpdateFields.SecondaryMentorUsername {
		t.Errorf("Project secondary mentor username did not get updated\n Expected: `%s`. Received: `%s`", projUpdateFields.SecondaryMentorUsername, updatedProj.SecondaryMentor.Username)
	}
}

// Test requests to /project/ with proper authentication and input
func TestProjectUpdateOK(t *testing.T) {
	// Set up a local test database path
	db := setTestDB()
	defer unsetTestDB(db)

	// Generate a jwt secret key for testing
	setTestJwtSecretKey()
	defer unsetTestJwtSecretKey()

	// Register a test mentor
	testUsername := getTestUsername()
	testLoginFields := utils.LoginJwtFields{
		Username: testUsername,
	}

	testJwt, _ := utils.GenerateLoginJwtString(testLoginFields)

	mentorReqFields := controllers.RegisterMentorReqFields{
		Username: testUsername,
		Email:    "testuser@example.com",
	}

	db.Table("mentors").Create(&models.Mentor{
		Username: mentorReqFields.Username,
		Email:    mentorReqFields.Email,
	})

	// Register a secondary test mentor
	mentorReqFields = controllers.RegisterMentorReqFields{
		Username: "testSecondary",
		Name:     "Secondary, Test",
		Email:    "testusersecond@example.com",
	}

	db.Table("mentors").Create(&models.Mentor{
		Username: mentorReqFields.Username,
		Name:     mentorReqFields.Name,
		Email:    mentorReqFields.Email,
	})

	// Non-existent project update test
	t.Run(
		"Test: non-existent project update.",
		func(t *testing.T) {
			tProjectUpdateNonExistent(db, testUsername, testJwt, t)
		},
	)

	// Existent project update test
	t.Run(
		"Test: Existent project update.",
		func(t *testing.T) {
			tProjectUpdateExistent(db, testUsername, testJwt, t)
		},
	)
}
