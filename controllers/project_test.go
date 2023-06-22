package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kwoc-backend/controllers"
	"kwoc-backend/utils"
	"math/rand"
	"net/http"
	"testing"

	"gorm.io/gorm"
)

func createProjctRegRequest(reqFields *controllers.RegisterProjectReqFields) *http.Request {
	reqBody, _ := json.Marshal(reqFields)

	req, _ := http.NewRequest(
		"POST",
		"/project/add/",
		bytes.NewReader(reqBody),
	)

	return req
}

func createTestProjectRegFields(mentorUsername string, secondaryMentorUsername string) *controllers.RegisterProjectReqFields {
	return &controllers.RegisterProjectReqFields{
		Name:                    fmt.Sprintf("YANGJF-%d", rand.Int()),
		Description:             "Yet another next-gen Javascript framework.",
		Tags:                    "next-gen, javascript, framework",
		MentorUsername:          mentorUsername,
		SecondaryMentorUsername: secondaryMentorUsername,
		RepoLink:                "https://xkcd.com/927/",
		ComChannel:              "com-channel",
		ReadmeURL:               "readme",
	}
}

// Test unauthenticated request to /project/add/
func TestProjectRegNoAuth(t *testing.T) {
	testRequestNoAuth(t, "POST", "/project/add/")
}

// Test request to /project/add/ with invalid jwt
func TestProjectRegInvalidAuth(t *testing.T) {
	testRequestInvalidAuth(t, "POST", "/project/add/")
}

// Test request to /project/add/ with session hijacking attempt
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
	expectResponseBodyToBe(t, res, "Login username and mentor username do not match.")
}

// Test a request to /project/add/ with non-existent mentors
func TestProjectRegInvalidMentor(t *testing.T) {
	// Set up a local test database path
	db := setTestDB()
	defer unsetTestDB()

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
	expectResponseBodyToBe(t, projectRes, fmt.Sprintf("Error: Mentor `%s` does not exist.", testUsername))
	// --- TEST PROJECT REGISTRATION WITH INVALID PRIMARY MENTOR ---
}

// Test a new project registration request to /project/add/ with proper authentication and input
func tProjectRegNew(db *gorm.DB, testUsername string, testJwt string, t *testing.T) {
	projectReqFields := createTestProjectRegFields(testUsername, "")

	projectReq := createProjctRegRequest(projectReqFields)
	projectReq.Header.Add("Bearer", testJwt)

	projectRes := executeRequest(projectReq, db)

	expectStatusCodeToBe(t, projectRes, http.StatusOK)
	expectResponseBodyToBe(t, projectRes, "Success.")
}

// Test an existing project registration request to /project/add/ with proper authentication and input
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
	expectResponseBodyToBe(t, projectRes, fmt.Sprintf("Error: Project `%s` already exists.", projectReqFields.RepoLink))
}

// Test requests to /project/add/ with proper authentication and input
func TestProjectRegOK(t *testing.T) {
	// Set up a local test database path
	db := setTestDB()
	defer unsetTestDB()

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
