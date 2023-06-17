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

// Test unauthenticated request to /student/form/
func TestStudentNoAuth(t *testing.T) {
	req, _ := http.NewRequest("POST", "/student/form/", nil)
	res := executeRequest(req)

	if res.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d. Got %d.", http.StatusUnauthorized, res.Code)
	}
}

// Test request to /student/form/ with invalid jwt
func TestStudentInvalidAuth(t *testing.T) {
	req, _ := http.NewRequest("POST", "/student/form/", nil)
	req.Header.Add("Bearer", "Some invalid token")

	res := executeRequest(req)

	// Expect internal server error because token parsing throws an error
	if res.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d. Got %d.", http.StatusInternalServerError, res.Code)
	}
}

// Test request to /student/form/ with session hijacking attempt
func TestStudentSessionHijacking(t *testing.T) {
	// Generate a jwt secret key for testing
	rand.Seed(time.Now().UnixMilli())

	os.Setenv("JWT_SECRET_KEY", fmt.Sprintf("testkey%d", rand.Int()))

	testLoginFields := utils.LoginJwtFields{
		Username: "someuser",
	}

	someuserJwt, _ := utils.GenerateLoginJwtString(testLoginFields)

	reqFields := controllers.RegisterStudentReqFields{
		Username: "anotheruser",
		Email:    "anotheruseremail@example.com",
		College:  "anotherusercollege",
	}

	reqBody, _ := json.Marshal(reqFields)

	req, _ := http.NewRequest(
		"POST",
		"/student/form/",
		bytes.NewReader(reqBody),
	)
	req.Header.Add("Bearer", someuserJwt)

	res := executeRequest(req)

	if res.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d. Got %d.", http.StatusInternalServerError, res.Code)
	}

	resBody := res.Body.String()
	expectedBody := "Login username and given username do not match."

	if resBody != expectedBody {
		t.Errorf("Expected response `%s`. Got `%s`.", expectedBody, resBody)
	}
}

// Test a request to /student/form/ with proper authentication and input
func TestStudentOK(t *testing.T) {
	// Set up a local test database path
	os.Setenv("DEV", "true")
	os.Setenv("DEV_DB_PATH", "testDB.db")
	utils.MigrateModels()

	// Generate a jwt secret key for testing
	rand.Seed(time.Now().UnixMilli())

	os.Setenv("JWT_SECRET_KEY", fmt.Sprintf("testkey%d", rand.Int()))

	// Test login fields
	testUsername := fmt.Sprintf("testuser%d", rand.Int())
	testLoginFields := utils.LoginJwtFields{
		Username: testUsername,
	}

	testJwt, _ := utils.GenerateLoginJwtString(testLoginFields)

	reqFields := controllers.RegisterStudentReqFields{
		Username: testUsername,
		Email:    "testuser@example.com",
		College:  "testusercollege",
	}

	reqBody, _ := json.Marshal(reqFields)

	// --- TEST NEW USER REGISTRATION ---
	req, _ := http.NewRequest(
		"POST",
		"/student/form/",
		bytes.NewReader(reqBody),
	)
	req.Header.Add("Bearer", testJwt)

	res := executeRequest(req)

	if res.Code != http.StatusOK {
		t.Errorf("Expected status code %d. Got %d.", http.StatusOK, res.Code)
	}

	resBody := res.Body.String()
	expectedBody := "Success."

	if resBody != expectedBody {
		t.Errorf("Expected response `%s`. Got `%s`.", expectedBody, resBody)
	}
	// --- TEST NEW USER REGISTRATION ---

	// --- TEST EXISTING USER REQUEST ---
	req, _ = http.NewRequest(
		"POST",
		"/student/form/",
		bytes.NewReader(reqBody),
	)
	req.Header.Add("Bearer", testJwt)

	res = executeRequest(req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d. Got %d.", http.StatusOK, res.Code)
	}

	resBody = res.Body.String()
	expectedBody = "Error: Student already exists."

	if resBody != expectedBody {
		t.Errorf("Expected response `%s`. Got `%s`.", expectedBody, resBody)
	}
	// --- TEST EXISTING USER REQUEST ---

	// Remove the test database
	os.Unsetenv("DEV_DB_PATH")
	os.Remove("testDB.db")
}
