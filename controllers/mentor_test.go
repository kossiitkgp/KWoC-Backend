package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kwoc-backend/controllers"
	"kwoc-backend/models"
	"kwoc-backend/utils"
	"net/http"
	"testing"

	"gorm.io/gorm"
)

func createMentorRegRequest(reqFields *controllers.RegisterMentorReqFields) *http.Request {
	reqBody, _ := json.Marshal(reqFields)

	req, _ := http.NewRequest(
		"POST",
		"/mentor/form/",
		bytes.NewReader(reqBody),
	)

	return req
}

// Test unauthenticated request to /mentor/form/
func TestMentorRegNoAuth(t *testing.T) {
	testRequestNoAuth(t, "POST", "/mentor/form/")
}

// Test request to /mentor/form/ with invalid jwt
func TestMentorRegInvalidAuth(t *testing.T) {
	testRequestInvalidAuth(t, "POST", "/mentor/form/")
}

// Test request to /mentor/form/ with session hijacking attempt
func TestMentorRegSessionHijacking(t *testing.T) {
	// Generate a jwt secret key for testing
	setTestJwtSecretKey()

	testLoginFields := utils.LoginJwtFields{Username: "someuser"}

	someuserJwt, _ := utils.GenerateLoginJwtString(testLoginFields)

	reqFields := controllers.RegisterMentorReqFields{Username: "anotheruser"}

	req := createMentorRegRequest(&reqFields)
	req.Header.Add("Bearer", someuserJwt)

	res := executeRequest(req, nil)

	expectStatusCodeToBe(t, res, http.StatusUnauthorized)
	expectResponseBodyToBe(t, res, "Login username and given username do not match.")
}

// Test a new user registration request to /mentor/form/ with proper authentication and input
func tMentorRegNewUser(db *gorm.DB, t *testing.T) {
	// Test login fields
	testUsername := getTestUsername()
	testLoginFields := utils.LoginJwtFields{Username: testUsername}

	testJwt, _ := utils.GenerateLoginJwtString(testLoginFields)
	reqFields := controllers.RegisterMentorReqFields{Username: testUsername}

	req := createMentorRegRequest(&reqFields)
	req.Header.Add("Bearer", testJwt)

	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusOK)
	expectResponseBodyToBe(t, res, "Mentor registration successful.")
}

// Test an existing user registration request to /mentor/form/ with proper authentication and input
func tMentorRegExistingUser(db *gorm.DB, t *testing.T) {
	// Test login fields
	testUsername := getTestUsername()
	testLoginFields := utils.LoginJwtFields{Username: testUsername}

	testJwt, _ := utils.GenerateLoginJwtString(testLoginFields)
	reqFields := controllers.RegisterMentorReqFields{Username: testUsername}

	req := createMentorRegRequest(&reqFields)
	req.Header.Add("Bearer", testJwt)

	_ = executeRequest(req, db)

	// Execute the same request again
	req = createMentorRegRequest(&reqFields)
	req.Header.Add("Bearer", testJwt)

	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusBadRequest)
	expectResponseBodyToBe(t, res, fmt.Sprintf("Mentor `%s` already exists.", testUsername))
}

// Test requests to /mentor/form/ with proper authentication and input
func TestMentorRegOK(t *testing.T) {
	// Set up a local test database path
	db := setTestDB()
	defer unsetTestDB()

	// Generate a jwt secret key for testing
	setTestJwtSecretKey()
	defer unsetTestJwtSecretKey()

	// New mentor registration test
	t.Run(
		"Test: new mentor registration.",
		func(t *testing.T) {
			tMentorRegNewUser(db, t)
		},
	)

	// Existing mentor registration test
	t.Run(
		"Test: existing mentor registration.",
		func(t *testing.T) {
			tMentorRegExistingUser(db, t)
		},
	)
}

func createFetchMentorRequest() *http.Request {
	req, _ := http.NewRequest(
		"GET",
		"/mentor/all/",
		nil,
	)
	return req
}

// Test unauthenticated request to /mentor/all/
func TestFetchMentorNoAuth(t *testing.T) {
	testRequestNoAuth(t, "GET", "/mentor/all/")
}

// Test request to /mentor/all/ with invalid jwt
func TestFetchMentorInvalidAuth(t *testing.T) {
	testRequestInvalidAuth(t, "GET", "/mentor/all/")
}

func TestFetchMentorOK(t *testing.T) {
	const numMentors = 10
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

	modelMentors := make([]models.Mentor, 0, numMentors)
	var testMentors [numMentors]controllers.Mentor
	for i := 0; i < numMentors; i++ {
		modelMentors = append(modelMentors,
			models.Mentor{
				Name:     fmt.Sprintf("Test%d", i),
				Username: fmt.Sprintf("test%d", i),
				Email:    fmt.Sprintf("test%d@example.com", i),
			})
		testMentors[i] = controllers.Mentor{
			Name:     fmt.Sprintf("Test%d", i),
			Username: fmt.Sprintf("test%d", i),
		}

	}
	_ = db.Table("mentors").Create(modelMentors)

	req := createFetchMentorRequest()
	req.Header.Add("Bearer", testJwt)

	res := executeRequest(req, db)

	var resMentors []controllers.Mentor
	_ = json.NewDecoder(res.Body).Decode(&resMentors)

	expectStatusCodeToBe(t, res, http.StatusOK)
	if len(resMentors) != 10 {
		t.Fatalf("Not getting expected numbers of mentors from /mentor/all/")
	}

	for i, mentor := range resMentors {
		if mentor != testMentors[i] {
			t.Fatalf("Incorrect mentors returned from /mentor/all/")
		}
	}
}
