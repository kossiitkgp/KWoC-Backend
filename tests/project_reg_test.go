package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"testing"

	"gorm.io/gorm"

	"github.com/kossiitkgp/kwoc-backend/v2/controllers"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"
)

func createTestProjectRegFields(mentorUsername string, secondaryMentorUsername string) *controllers.RegisterProjectReqFields {
	return &controllers.RegisterProjectReqFields{
		Name:                    fmt.Sprintf("YANGJF-%d", rand.Int()),
		Description:             "Yet another next-gen Javascript framework.",
		Tags:                    strings.Split("next-gen,javascript,framework", ","),
		MentorUsername:          mentorUsername,
		SecondaryMentorUsername: secondaryMentorUsername,
		RepoLink:                "https://xkcd.com/927/",
		CommChannel:             "comm-channel",
		ReadmeLink:              "readme",
	}
}

func createProjctRegRequest(reqFields *controllers.RegisterProjectReqFields) *http.Request {
	reqBody, _ := json.Marshal(reqFields)

	req, _ := http.NewRequest(
		"POST",
		"/project/",
		bytes.NewReader(reqBody),
	)

	return req
}

// Test unauthenticated request to /project/
func TestProjectRegNoAuth(t *testing.T) {
	testRequestNoAuth(t, "POST", "/project/")
}

// Test request to /project/ with invalid jwt
func TestProjectRegInvalidAuth(t *testing.T) {
	testRequestInvalidAuth(t, "POST", "/project/")
}

// Test request to /project/ with session hijacking attempt
func TestProjectRegSessionHijacking(t *testing.T) {
	// Generate a jwt secret key for testing
	setTestJwtSecretKey()
	defer unsetTestJwtSecretKey()

	testLoginFields := utils.LoginJwtFields{Username: "someuser"}

	someuserJwt, _ := utils.GenerateLoginJwtString(testLoginFields)

	reqFields := controllers.RegisterProjectReqFields{MentorUsername: "anotheruser"}

	req := createProjctRegRequest(&reqFields)
	req.Header.Add("Bearer", someuserJwt)

	res := executeRequest(req, nil)

	expectStatusCodeToBe(t, res, http.StatusUnauthorized)
	expectResponseJSONBodyToBe(t, res, utils.HTTPMessage{StatusCode: http.StatusUnauthorized, Message: "Login username and mentor username do not match."})
}

// Test a request to /project/ with non-existent mentors
func TestProjectRegInvalidMentor(t *testing.T) {
	// Set up a local test database path
	db := setTestDB()
	defer unsetTestDB(db)

	// Generate a jwt secret key for testing
	setTestJwtSecretKey()
	defer unsetTestJwtSecretKey()

	// Register a test mentor
	testUsername := getTestUsername()
	testLoginFields := utils.LoginJwtFields{Username: testUsername}

	testJwt, _ := utils.GenerateLoginJwtString(testLoginFields)

	// --- TEST PROJECT REGISTRATION WITH INVALID PRIMARY MENTOR ---
	projectReqFields := createTestProjectRegFields(testUsername, "")

	projectReq := createProjctRegRequest(projectReqFields)
	projectReq.Header.Add("Bearer", testJwt)

	projectRes := executeRequest(projectReq, db)

	expectStatusCodeToBe(t, projectRes, http.StatusBadRequest)
	expectResponseJSONBodyToBe(t, projectRes, utils.HTTPMessage{StatusCode: http.StatusBadRequest, Message: fmt.Sprintf("Error: Mentor `%s` does not exist.", testUsername)})
	// --- TEST PROJECT REGISTRATION WITH INVALID PRIMARY MENTOR ---
}

// Test a new project registration request to /project/ with proper authentication and input
func tProjectRegNew(db *gorm.DB, testUsername string, testJwt string, t *testing.T) {
	projectReqFields := createTestProjectRegFields(testUsername, "")

	projectReq := createProjctRegRequest(projectReqFields)
	projectReq.Header.Add("Bearer", testJwt)

	projectRes := executeRequest(projectReq, db)

	expectStatusCodeToBe(t, projectRes, http.StatusOK)
	expectResponseJSONBodyToBe(t, projectRes, utils.HTTPMessage{StatusCode: http.StatusOK, Message: "Success."})
}

// Test an existing project registration request to /project/ with proper authentication and input
func tProjectRegExisting(db *gorm.DB, testUsername string, testJwt string, t *testing.T) {
	projectReqFields := createTestProjectRegFields(testUsername, "")

	projectReq := createProjctRegRequest(projectReqFields)
	projectReq.Header.Add("Bearer", testJwt)

	_ = executeRequest(projectReq, db)

	// Execute the same request again
	projectReq = createProjctRegRequest(projectReqFields)
	projectReq.Header.Add("Bearer", testJwt)

	projectRes := executeRequest(projectReq, db)

	expectStatusCodeToBe(t, projectRes, http.StatusBadRequest)
	expectResponseJSONBodyToBe(t, projectRes, utils.HTTPMessage{StatusCode: http.StatusBadRequest, Message: fmt.Sprintf("Error: Project `%s` already exists.", projectReqFields.RepoLink)})
}

// Test requests to /project/ with proper authentication and input
func TestProjectRegOK(t *testing.T) {
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

	mentorReq := createMentorRegRequest(&mentorReqFields)
	mentorReq.Header.Add("Bearer", testJwt)
	_ = executeRequest(mentorReq, db)

	// New project registration test
	t.Run(
		"Test: new project registration.",
		func(t *testing.T) {
			tProjectRegNew(db, testUsername, testJwt, t)
		},
	)

	// Existing project registration test
	t.Run(
		"Test: existing project registration.",
		func(t *testing.T) {
			tProjectRegExisting(db, testUsername, testJwt, t)
		},
	)
}
