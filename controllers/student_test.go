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
	testRequestNoAuth(t, "POST", "/student/form/")
}

// Test request to /student/form/ with invalid jwt
func TestStudentInvalidAuth(t *testing.T) {
	testRequestInvalidAuth(t, "POST", "/student/form/")
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

	res := executeRequest(req, nil)

	expectStatusCodeToBe(t, res, http.StatusUnauthorized)
	expectResponseBodyToBe(t, res, "Login username and given username do not match.")
}

// Test a request to /student/form/ with proper authentication and input
func TestStudentOK(t *testing.T) {
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

	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusOK)
	expectResponseBodyToBe(t, res, "Student registration successful.")
	// --- TEST NEW USER REGISTRATION ---

	// --- TEST EXISTING USER REQUEST ---
	req, _ = http.NewRequest(
		"POST",
		"/student/form/",
		bytes.NewReader(reqBody),
	)
	req.Header.Add("Bearer", testJwt)

	res = executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusBadRequest)
	expectResponseBodyToBe(t, res, fmt.Sprintf("Student `%s` already exists.", testUsername))
	// --- TEST EXISTING USER REQUEST ---
}
