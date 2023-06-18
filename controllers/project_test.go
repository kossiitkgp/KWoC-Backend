package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kwoc-backend/controllers"
	"kwoc-backend/utils"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"
)

// Test unauthenticated request to /project/add/
func TestProjectNoAuth(t *testing.T) {
	testRequestNoAuth(t, "POST", "/project/add/")
}

// Test request to /project/add/ with invalid jwt
func TestProjectInvalidAuth(t *testing.T) {
	testRequestInvalidAuth(t, "POST", "/project/add/")
}

// Test request to /project/add/ with session hijacking attempt
func TestProjectSessionHijacking(t *testing.T) {
	// Generate a jwt secret key for testing
	rand.Seed(time.Now().UnixMilli())

	os.Setenv("JWT_SECRET_KEY", fmt.Sprintf("testkey%d", rand.Int()))

	testLoginFields := utils.LoginJwtFields{
		Username: "someuser",
	}

	someuserJwt, _ := utils.GenerateLoginJwtString(testLoginFields)

	reqFields := controllers.RegisterProjectReqFields{
		MentorUsername: "anotheruser",
	}

	reqBody, _ := json.Marshal(reqFields)

	req, _ := http.NewRequest(
		"POST",
		"/project/add/",
		bytes.NewReader(reqBody),
	)
	req.Header.Add("Bearer", someuserJwt)

	res := executeRequest(req, nil)

	expectStatusCodeToBe(t, res, http.StatusUnauthorized)
	expectResponseBodyToBe(t, res, "Login username and mentor username do not match.")
}

// Test a request to /project/add/ with non-existent mentors
func TestProjectInvalidMentor(t *testing.T) {
	// Set up a local test database path
	os.Setenv("DEV", "true")
	os.Setenv("DEV_DB_PATH", "testDB.db")
	db, _ := utils.GetDB()
	_ = utils.MigrateModels(db)
	// Remove the test database once used
	defer os.Unsetenv("DEV_DB_PATH")
	defer os.Remove("testDB.db")

	// Generate a jwt secret key for testing
	rand.Seed(time.Now().UnixMilli())

	os.Setenv("JWT_SECRET_KEY", fmt.Sprintf("testkey%d", rand.Int()))

	// Register a test mentor
	testUsername := fmt.Sprintf("testuser%d", rand.Int())
	testLoginFields := utils.LoginJwtFields{
		Username: testUsername,
	}

	testJwt, _ := utils.GenerateLoginJwtString(testLoginFields)

	// --- TEST PROJECT REGISTRATION WITH INVALID PRIMARY MENTOR ---
	projectReqFields := controllers.RegisterProjectReqFields{
		Name:                    "YANGJF",
		Description:             "Yet another next-gen Javascript framework.",
		Tags:                    "next-gen, javascript, framework",
		MentorUsername:          testUsername,
		SecondaryMentorUsername: "",
		RepoLink:                "https://xkcd.com/927/",
		ComChannel:              "com-channel",
		ReadmeURL:               "readme",
	}

	projectReqBody, _ := json.Marshal(projectReqFields)

	projectReq, _ := http.NewRequest(
		"POST",
		"/project/add/",
		bytes.NewReader(projectReqBody),
	)
	projectReq.Header.Add("Bearer", testJwt)

	projectRes := executeRequest(projectReq, db)

	expectStatusCodeToBe(t, projectRes, http.StatusBadRequest)
	expectResponseBodyToBe(t, projectRes, fmt.Sprintf("Error: Mentor `%s` does not exist.", testUsername))
	// --- TEST PROJECT REGISTRATION WITH INVALID PRIMARY MENTOR ---
}

// Test a request to /project/add/ with proper authentication and input
func TestProjectOK(t *testing.T) {
	// Set up a local test database path
	os.Setenv("DEV", "true")
	os.Setenv("DEV_DB_PATH", "testDB.db")
	db, _ := utils.GetDB()
	_ = utils.MigrateModels(db)
	// Remove the test database once used
	defer os.Unsetenv("DEV_DB_PATH")
	defer os.Remove("testDB.db")

	// Generate a jwt secret key for testing
	rand.Seed(time.Now().UnixMilli())

	os.Setenv("JWT_SECRET_KEY", fmt.Sprintf("testkey%d", rand.Int()))

	// Register a test mentor
	testUsername := fmt.Sprintf("testuser%d", rand.Int())
	testLoginFields := utils.LoginJwtFields{
		Username: testUsername,
	}

	testJwt, _ := utils.GenerateLoginJwtString(testLoginFields)

	mentorReqFields := controllers.RegisterMentorReqFields{
		Username: testUsername,
		Email:    "testuser@example.com",
	}
	mentorReqBody, _ := json.Marshal(mentorReqFields)

	mentorReq, _ := http.NewRequest(
		"POST",
		"/mentor/form/",
		bytes.NewReader(mentorReqBody),
	)
	mentorReq.Header.Add("Bearer", testJwt)
	_ = executeRequest(mentorReq, db)

	// --- TEST NEW PROJECT REGISTRATION ---
	projectReqFields := controllers.RegisterProjectReqFields{
		Name:                    "YANGJF",
		Description:             "Yet another next-gen Javascript framework.",
		Tags:                    "next-gen, javascript, framework",
		MentorUsername:          testUsername,
		SecondaryMentorUsername: "",
		RepoLink:                "https://xkcd.com/927/",
		ComChannel:              "com-channel",
		ReadmeURL:               "readme",
	}

	projectReqBody, _ := json.Marshal(projectReqFields)

	projectReq, _ := http.NewRequest(
		"POST",
		"/project/add/",
		bytes.NewReader(projectReqBody),
	)
	projectReq.Header.Add("Bearer", testJwt)

	projectRes := executeRequest(projectReq, db)

	expectStatusCodeToBe(t, projectRes, http.StatusOK)
	expectResponseBodyToBe(t, projectRes, "Success.")
	// --- TEST NEW PROJECT REGISTRATION ---

	// --- TEST EXISTING PROJECT REGISTRATION ---
	projectReq, _ = http.NewRequest(
		"POST",
		"/project/add/",
		bytes.NewReader(projectReqBody),
	)
	projectReq.Header.Add("Bearer", testJwt)

	projectRes = executeRequest(projectReq, db)

	expectStatusCodeToBe(t, projectRes, http.StatusBadRequest)
	expectResponseBodyToBe(t, projectRes, fmt.Sprintf("Error: Project `%s` already exists.", projectReqFields.RepoLink))
	// --- TEST EXISTING PROJECT REGISTRATION ---
}
