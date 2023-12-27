package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"gorm.io/gorm"

	"github.com/kossiitkgp/kwoc-backend/v2/controllers"
	"github.com/kossiitkgp/kwoc-backend/v2/models"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"
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
	expectResponseJSONBodyToBe(t, res, utils.HTTPMessage{StatusCode: http.StatusUnauthorized, Message: "Login username and given username do not match."})
}

// Test unauthenticated request to /mentor/form/ [put]
func TestMentorUpdateNoAuth(t *testing.T) {
	testRequestNoAuth(t, "PUT", "/mentor/form/")
}

// Test request to /mentor/form/ [put] with invalid jwt
func TestMentorUpdateInvalidAuth(t *testing.T) {
	testRequestInvalidAuth(t, "PUT", "/mentor/form/")
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
	expectResponseJSONBodyToBe(t, res, utils.HTTPMessage{StatusCode: http.StatusOK, Message: "Mentor registration successful."})
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
	expectResponseJSONBodyToBe(t, res, utils.HTTPMessage{StatusCode: http.StatusBadRequest, Message: fmt.Sprintf("Mentor `%s` already exists.", testUsername)})
}

// Test an existing student registration request to /mentor/form/ with proper authentication and input
func tMentorRegAsStudent(db *gorm.DB, t *testing.T) {
	// Test login fields
	testUsername := getTestUsername()
	testLoginFields := utils.LoginJwtFields{Username: testUsername}

	testJwt, _ := utils.GenerateLoginJwtString(testLoginFields)

	db.Table("students").Create(&models.Student{Username: testUsername})

	mentorFields := controllers.RegisterMentorReqFields{Username: testUsername}
	req := createMentorRegRequest(&mentorFields)
	req.Header.Add("Bearer", testJwt)

	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusBadRequest)
	expectResponseJSONBodyToBe(t, res, utils.HTTPMessage{StatusCode: http.StatusBadRequest, Message: fmt.Sprintf("The username `%s` already exists as a student.", testUsername)})
}

// Test requests to /mentor/form/ with proper authentication and input
func TestMentorRegOK(t *testing.T) {
	// Set up a local test database path
	db := setTestDB()
	defer unsetTestDB(db)

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

	// Student registering as mentor test
	t.Run(
		"Test: Student registering as mentor.",
		func(t *testing.T) {
			tMentorRegAsStudent(db, t)
		},
	)
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
	defer unsetTestDB(db)

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
	expectResponseJSONBodyToBe(t, res, utils.HTTPMessage{StatusCode: http.StatusBadRequest, Message: fmt.Sprintf("Mentor `%s` does not exists.", testUsername)})
}

// Test requests to /mentor/dashboard/ with registered and  proper authentication
func TestMentorDashboardOK(t *testing.T) {
	// Set up a local test database path
	db := setTestDB()
	defer unsetTestDB(db)

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

	db.Table("mentors").Create(&modelMentor)

	mentorID := int32(modelMentor.ID)
	testProjects := generateTestProjects(5, false, true)
	for i := range testProjects {
		testProjects[i].MentorId = mentorID
		testProjects[i].SecondaryMentorId = &mentorID
	}
	var projects []controllers.ProjectInfo
	var students []controllers.StudentInfo

	modelStudents := generateTestStudents(5)

	for i, student := range modelStudents {
		if i < 3 {
			testProjects[1].Contributors = testProjects[1].Contributors + student.Username + ","
		} else {
			testProjects[3].Contributors = testProjects[3].Contributors + student.Username + ","
		}
	}

	testProjects[1].Contributors = strings.TrimSuffix(testProjects[1].Contributors, ",")
	testProjects[3].Contributors = strings.TrimSuffix(testProjects[3].Contributors, ",")
	db.Table("projects").Create(testProjects)

	for _, p := range testProjects {
		if (p.MentorId != int32(modelMentor.ID)) && (p.SecondaryMentorId != &mentorID) {
			continue
		}

		if p.MentorId == int32(modelMentor.ID) {
			p.Mentor = models.Mentor{
				Name:     modelMentor.Name,
				Username: modelMentor.Username,
			}
		}
		if p.SecondaryMentorId == &mentorID {
			p.SecondaryMentor = models.Mentor{
				Name:     modelMentor.Name,
				Username: modelMentor.Username,
			}
		}

		pulls := make([]string, 0)
		if len(p.Pulls) > 0 {
			pulls = strings.Split(p.Pulls, ",")
		}
		tags := make([]string, 0)
		if len(p.Tags) > 0 {
			tags = strings.Split(p.Tags, ",")
		}
		projects = append(projects, controllers.ProjectInfo{
			Id:            p.ID,
			Name:          p.Name,
			Description:   p.Description,
			RepoLink:      p.RepoLink,
			ReadmeLink:    p.ReadmeLink,
			Tags:          tags,
			ProjectStatus: p.ProjectStatus,
			StatusRemark:  p.StatusRemark,

			CommitCount:  p.CommitCount,
			PullCount:    p.PullCount,
			LinesAdded:   p.LinesAdded,
			LinesRemoved: p.LinesRemoved,

			Pulls: pulls,
			Mentor: controllers.Mentor{
				Username: p.Mentor.Username,
				Name:     p.Mentor.Name,
			},
			SecondaryMentor: controllers.Mentor{
				Username: p.SecondaryMentor.Username,
				Name:     p.SecondaryMentor.Name,
			},
		})
	}

	for _, student := range modelStudents {
		students = append(students, controllers.StudentInfo{
			Name:     student.Name,
			Username: student.Username,
		})
	}

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

	expectStatusCodeToBe(t, res, http.StatusOK)
	if !reflect.DeepEqual(testMentor, resMentor) {
		fmt.Printf("Expected mentor dashboard: %#v\n\n", testMentor)
		fmt.Printf("Received mentor dashboard: %#v\n", resMentor)
		t.Fatalf("Incorrect data returned from /mentor/dashboard/")
	}
}
