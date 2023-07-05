package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kwoc-backend/controllers"
	"kwoc-backend/models"
	"kwoc-backend/utils"
	"math/rand"
	"net/http"
	"strings"
	"testing"

	"gorm.io/gorm"
)

func createStudentRegRequest(reqFields *controllers.RegisterStudentReqFields) *http.Request {
	reqBody, _ := json.Marshal(reqFields)

	req, _ := http.NewRequest(
		"POST",
		"/student/form/",
		bytes.NewReader(reqBody),
	)

	return req
}

// Test unauthenticated request to /student/form/
func TestStudentRegNoAuth(t *testing.T) {
	testRequestNoAuth(t, "POST", "/student/form/")
}

// Test request to /student/form/ with invalid jwt
func TestStudentRegInvalidAuth(t *testing.T) {
	testRequestInvalidAuth(t, "POST", "/student/form/")
}

// Test request to /student/form/ with session hijacking attempt
func TestStudentRegSessionHijacking(t *testing.T) {
	// Generate a jwt secret key for testing
	setTestJwtSecretKey()

	testLoginFields := utils.LoginJwtFields{Username: "someuser"}

	someuserJwt, _ := utils.GenerateLoginJwtString(testLoginFields)

	reqFields := controllers.RegisterStudentReqFields{Username: "anotheruser"}

	req := createStudentRegRequest(&reqFields)
	req.Header.Add("Bearer", someuserJwt)

	res := executeRequest(req, nil)

	expectStatusCodeToBe(t, res, http.StatusUnauthorized)
	expectResponseBodyToBe(t, res, "Login username and given username do not match.")
}

// Test a new user registration request to /student/form/ with proper authentication and input
func tStudentRegNewUser(db *gorm.DB, t *testing.T) {
	// Test login fields
	testUsername := getTestUsername()
	testLoginFields := utils.LoginJwtFields{Username: testUsername}

	testJwt, _ := utils.GenerateLoginJwtString(testLoginFields)
	reqFields := controllers.RegisterStudentReqFields{Username: testUsername}

	req := createStudentRegRequest(&reqFields)
	req.Header.Add("Bearer", testJwt)

	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusOK)
	expectResponseBodyToBe(t, res, "Student registration successful.")
}

// Test an existing user registration request to /student/form/ with proper authentication and input
func tStudentRegExistingUser(db *gorm.DB, t *testing.T) {
	// Test login fields
	testUsername := getTestUsername()
	testLoginFields := utils.LoginJwtFields{Username: testUsername}

	testJwt, _ := utils.GenerateLoginJwtString(testLoginFields)
	reqFields := controllers.RegisterStudentReqFields{Username: testUsername}

	req := createStudentRegRequest(&reqFields)
	req.Header.Add("Bearer", testJwt)

	_ = executeRequest(req, db)

	// Execute the same request again
	req = createStudentRegRequest(&reqFields)
	req.Header.Add("Bearer", testJwt)

	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusBadRequest)
	expectResponseBodyToBe(t, res, fmt.Sprintf("Student `%s` already exists.", testUsername))
}

// Test requests to /student/form/ with proper authentication and input
func TestStudentRegOK(t *testing.T) {
	// Set up a local test database path
	db := setTestDB()
	defer unsetTestDB()

	// Generate a jwt secret key for testing
	setTestJwtSecretKey()
	defer unsetTestJwtSecretKey()

	// New student registration test
	t.Run(
		"Test: new student registration.",
		func(t *testing.T) {
			tStudentRegNewUser(db, t)
		},
	)

	// Existing student registration test
	t.Run(
		"Test: existing student registration.",
		func(t *testing.T) {
			tStudentRegExistingUser(db, t)
		},
	)
}

func createStudentBlogLinkRequest(reqFields *controllers.StudentBlogLinkReqFields) *http.Request {
	reqBody, _ := json.Marshal(reqFields)

	req, _ := http.NewRequest(
		"POST",
		"/student/bloglink/",
		bytes.NewReader(reqBody),
	)

	return req
}

// Test unauthenticated request to /student/bloglink/
func TestStudentBloglinkNoAuth(t *testing.T) {
	testRequestNoAuth(t, "POST", "/student/bloglink/")
}

// Test request to /student/bloglink/ with invalid jwt
func TestStudentBloglinkInvalidAuth(t *testing.T) {
	testRequestInvalidAuth(t, "POST", "/student/bloglink/")
}

// Test request to /student/bloglink/ with session hijacking attempt
func TestStudentBloglinkSessionHijacking(t *testing.T) {
	// Generate a jwt secret key for testing
	setTestJwtSecretKey()
	defer unsetTestJwtSecretKey()

	testLoginFields := utils.LoginJwtFields{Username: "someuser"}

	someuserJwt, _ := utils.GenerateLoginJwtString(testLoginFields)

	reqFields := controllers.StudentBlogLinkReqFields{Username: "anotheruser", BlogLink: "https://grugbrain.dev"}

	req := createStudentBlogLinkRequest(&reqFields)
	req.Header.Add("Bearer", someuserJwt)

	res := executeRequest(req, nil)

	expectStatusCodeToBe(t, res, http.StatusUnauthorized)
	expectResponseBodyToBe(t, res, "Login username and given username do not match.")
}

// Test an existing user request to /student/bloglink/ with proper authentication and input
func tStudentBlogLinkExistingUser(db *gorm.DB, t *testing.T) {
	// Test login fields
	testUsername := getTestUsername()
	testLoginFields := utils.LoginJwtFields{Username: testUsername}

	testJwt, _ := utils.GenerateLoginJwtString(testLoginFields)
	reqFieldsReg := controllers.RegisterStudentReqFields{Username: testUsername}

	req := createStudentRegRequest(&reqFieldsReg)
	req.Header.Add("Bearer", testJwt)

	_ = executeRequest(req, db)

	// Execute the bloglink request
	reqFields := controllers.StudentBlogLinkReqFields{Username: testUsername, BlogLink: "https://grugbrain.dev/"}
	req = createStudentBlogLinkRequest(&reqFields)
	req.Header.Add("Bearer", testJwt)

	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusOK)
	expectResponseBodyToBe(t, res, "BlogLink successfully updated.")
}

// Test a non existing/registered user's request to /student/bloglink/ with proper authentication and input
func tStudentBlogLinkNonExistingUser(db *gorm.DB, t *testing.T) {
	// Test login fields
	testUsername := getTestUsername()
	testLoginFields := utils.LoginJwtFields{Username: testUsername}

	testJwt, _ := utils.GenerateLoginJwtString(testLoginFields)

	// Execute the bloglink request
	reqFields := controllers.StudentBlogLinkReqFields{Username: testUsername, BlogLink: "https://grugbrain.dev/"}
	req := createStudentBlogLinkRequest(&reqFields)
	req.Header.Add("Bearer", testJwt)

	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusBadRequest)
	expectResponseBodyToBe(t, res, fmt.Sprintf("Student `%s` does not exists.", testUsername))
}

// Test request  /student/bloglink/ with proper authentication and input
func TestStudentBlogLink(t *testing.T) {
	// Set up a local test database path
	db := setTestDB()
	defer unsetTestDB()

	// Generate a jwt secret key for testing
	setTestJwtSecretKey()
	defer unsetTestJwtSecretKey()

	// Existing student test
	t.Run(
		"Test: existing student bloglink request",
		func(t *testing.T) {
			tStudentBlogLinkExistingUser(db, t)
		},
	)

	// Non Existing student test
	t.Run(
		"Test: non existing student bloglink request",
		func(t *testing.T) {
			tStudentBlogLinkNonExistingUser(db, t)
		},
	)
}

// Test unauthenticated request to /student/dashboard/
func TestStudentDashboardNoAuth(t *testing.T) {
	testRequestNoAuth(t, "GET", "/student/dashboard/")
}

// Test request to /student/dashboard/ with invalid jwt
func TestStudentDashboardInvalidAuth(t *testing.T) {
	testRequestInvalidAuth(t, "GET", "/student/dashboard/")
}

// Test requests to /student/dashboard/ with proper authentication
func TestStudentDashboardOK(t *testing.T) {
	// Set up a local test database path
	db := setTestDB()
	defer unsetTestDB()

	// Generate a jwt secret key for testing
	setTestJwtSecretKey()
	defer unsetTestJwtSecretKey()

	// Test login fields
	testUsername := getTestUsername()
	testLoginFields := utils.LoginJwtFields{Username: testUsername}

	testJwt, _ := utils.GenerateLoginJwtString(testLoginFields)

	testProjects := generateTestProjects(5, false, true)
	_ = db.Table("projects").Create(testProjects)

	var project_ids []string
	for _, p := range testProjects {
		project_ids = append(project_ids, fmt.Sprint(p.ID))

	}
	modelStudent := models.Student{
		Name:           "Test",
		Username:       testUsername,
		College:        "The best university",
		ProjectsWorked: strings.Join(project_ids, ","),
		PassedMidEvals: true,
		PassedEndEvals: true,
		CommitCount:    uint(rand.Int()),
		PullCount:      uint(rand.Int()),
		LinesAdded:     uint(rand.Int()),
		LinesRemoved:   uint(rand.Int()),
	}

	_ = db.Table("students").Create(&modelStudent)
	testStudent := controllers.CreateStudentDashboard(modelStudent, db)

	req, _ := http.NewRequest(
		"GET",
		"/student/dashboard/",
		nil,
	)
	req.Header.Add("Bearer", testJwt)

	res := executeRequest(req, db)

	var resStudent controllers.StudentDashboard
	_ = json.NewDecoder(res.Body).Decode(&resStudent)

	expectStatusCodeToBe(t, res, http.StatusOK)
	if fmt.Sprint(testStudent) != fmt.Sprint(resStudent) {
		t.Fatalf("Incorrect data returned from /student/dashboard/")
	}
}
