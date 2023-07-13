package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/kossiitkgp/kwoc-backend/v2/controllers"
	"github.com/kossiitkgp/kwoc-backend/v2/models"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"

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
	if len(resMentors) != numMentors {
		t.Fatalf("Not getting expected numbers of mentors from /mentor/all/")
	}

	for i, mentor := range resMentors {
		if mentor != testMentors[i] {
			t.Fatalf("Incorrect mentors returned from /mentor/all/")
		}
	}
}

// Test unauthenticated request to /mentor/dashboard/
func TestMentorDashboardNoAuth(t *testing.T) {
	testRequestNoAuth(t, "GET", "/mentor/dashboard/")
}

// Test request to /mentor/dashboard/ with invalid jwt
func TestMentorDashboardInvalidAuth(t *testing.T) {
	testRequestInvalidAuth(t, "GET", "/mentor/dashboard/")
}

// Test unauthenticated request to /mentor/dashboard/ with no registration
func TestMentorDashboardNoReg(t *testing.T) {
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

	req, _ := http.NewRequest(
		"GET",
		"/mentor/dashboard/",
		nil,
	)
	req.Header.Add("Bearer", testJwt)

	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusBadRequest)
	expectResponseBodyToBe(t, res, fmt.Sprintf("Mentor `%s` does not exists.", testUsername))
}

// Test requests to /mentor/dashboard/ with registered and  proper authentication
func TestMentorDashboardOK(t *testing.T) {
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

	modelMentor := models.Mentor{
		Name:     "TestMentor",
		Email:    "iamamentor@cool.com",
		Username: testUsername,
	}

	testProjects := generateTestProjects(5, false, true)
	testProjects[1].Mentor = modelMentor
	testProjects[1].MentorId = int32(modelMentor.ID)
	testProjects[3].SecondaryMentor = modelMentor
	testProjects[3].SecondaryMentorId = int32(modelMentor.ID)

	db.Table("mentors").Create(&modelMentor)

	var projects []controllers.ProjectInfo
	var students []controllers.StudentDashboard

	modelStudents := generateTestStudents(5)

	for i, student := range modelStudents {
		if i < 3 {
			student.ProjectsWorked = fmt.Sprint(testProjects[1].ID)
			testProjects[1].Contributors = testProjects[1].Contributors + "," + student.Username
		} else {
			student.ProjectsWorked = fmt.Sprint(testProjects[3].ID)
			testProjects[3].Contributors = testProjects[3].Contributors + "," + student.Username
		}
	}

	for _, p := range testProjects {
		if p.MentorId != int32(modelMentor.ID) && p.SecondaryMentorId != int32(modelMentor.ID) {
			continue
		}
		p.RepoLink = "www.thisisaLink.com"

		projects = append(projects, controllers.ProjectInfo{
			Name:     p.Name,
			RepoLink: p.RepoLink,

			CommitCount:  p.CommitCount,
			PullCount:    p.PullCount,
			LinesAdded:   p.LinesAdded,
			LinesRemoved: p.LinesRemoved,
		})
	}

	for _, student := range modelStudents {
		students = append(students, controllers.CreateStudentDashboard(student, db))
	}

	db.Table("projects").Create(testProjects)
	db.Table("students").Create(modelStudents)

	testMentor := controllers.MentorDashboard{
		Name:     modelMentor.Name,
		Username: modelMentor.Username,
		Email:    modelMentor.Email,

		Projects: projects,
		Students: students,
	}

	req, _ := http.NewRequest(
		"GET",
		"/mentor/dashboard/",
		nil,
	)
	req.Header.Add("Bearer", testJwt)

	res := executeRequest(req, db)

	var resMentor controllers.MentorDashboard
	_ = json.NewDecoder(res.Body).Decode(&resMentor)

	studentsEqual := true
	for _, student := range resMentor.Students {
		for _, modStudent := range modelStudents {
			if modStudent.Username != student.Username {
				studentsEqual = false
			}
		}
	}

	projectsEqual := true
	for _, proj := range resMentor.Projects {
		if !(proj.Name == projects[1].Name || proj.Name == projects[3].Name) {
			projectsEqual = false
		}
	}

	expectStatusCodeToBe(t, res, http.StatusOK)
	if !(testMentor.Name == resMentor.Name &&
		testMentor.Email == resMentor.Email &&
		testMentor.Username == resMentor.Username &&
		projectsEqual &&
		studentsEqual) {
		t.Fatalf("Incorrect data returned from /mentor/dashboard/")
	}
}
