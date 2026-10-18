package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kossiitkgp/kwoc-backend/v2/controllers"
	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
	"github.com/kossiitkgp/kwoc-backend/v2/models"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"
)

func createVerifyEmailRequest(reqFields *controllers.VerifyEmailReqFields) *http.Request {
	reqBody, _ := json.Marshal(reqFields)

	req, _ := http.NewRequest(
		"POST",
		"/verify/email/",
		bytes.NewReader(reqBody),
	)

	return req
}

func TestVerifyEmailNoParams(t *testing.T) {
	db := setTestDB()
	defer unsetTestDB(db)

	reqFields := controllers.VerifyEmailReqFields{}
	req := createVerifyEmailRequest(&reqFields)

	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusBadRequest)
	expectResponseJSONBodyToBe(t, res, utils.HTTPMessage{
		StatusCode: http.StatusBadRequest,
		Message:    "Provide a verification token, or a username and code.",
	})
}

func TestVerifyEmailInvalidToken(t *testing.T) {
	db := setTestDB()
	defer unsetTestDB(db)

	reqFields := controllers.VerifyEmailReqFields{Token: "not-a-real-token"}
	req := createVerifyEmailRequest(&reqFields)

	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusBadRequest)
	expectResponseJSONBodyToBe(t, res, utils.HTTPMessage{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid or expired verification code.",
	})
}

func TestVerifyEmailWithToken(t *testing.T) {
	db := setTestDB()
	defer unsetTestDB(db)

	testUsername := getTestUsername()
	student := models.Student{
		Username: testUsername,
		Name:     "Test Student",
		Email:    "student@example.com",
	}
	if err := db.Create(&student).Error; err != nil {
		t.Fatalf("Error creating test student: %v", err)
	}

	challenge, err := utils.CreateVerificationChallenge(db, testUsername, student.Email)
	if err != nil {
		t.Fatalf("Error creating verification challenge: %v", err)
	}

	reqFields := controllers.VerifyEmailReqFields{Token: challenge.Token}
	req := createVerifyEmailRequest(&reqFields)

	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusOK)

	var resFields controllers.VerifyEmailResFields
	if err := json.NewDecoder(res.Body).Decode(&resFields); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if resFields.Username != testUsername || resFields.Role != "student" || !resFields.EmailVerified {
		t.Errorf("Unexpected verification response: %+v", resFields)
	}

	updatedStudent := models.Student{}
	if err := db.Where("username = ?", testUsername).First(&updatedStudent).Error; err != nil {
		t.Fatalf("Error fetching student: %v", err)
	}

	if !updatedStudent.EmailVerified {
		t.Error("Expected student to be marked as verified.")
	}

	// Replaying the same token must fail.
	req = createVerifyEmailRequest(&reqFields)
	res = executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusBadRequest)
	expectResponseJSONBodyToBe(t, res, utils.HTTPMessage{
		StatusCode: http.StatusBadRequest,
		Message:    "This email address has already been verified.",
	})
}

func TestVerifyEmailWithCode(t *testing.T) {
	db := setTestDB()
	defer unsetTestDB(db)

	testUsername := getTestUsername()
	mentor := models.Mentor{
		Username: testUsername,
		Name:     "Test Mentor",
		Email:    "mentor@example.com",
	}
	if err := db.Create(&mentor).Error; err != nil {
		t.Fatalf("Error creating test mentor: %v", err)
	}

	challenge, err := utils.CreateVerificationChallenge(db, testUsername, mentor.Email)
	if err != nil {
		t.Fatalf("Error creating verification challenge: %v", err)
	}

	// An incorrect code should not verify the account.
	reqFields := controllers.VerifyEmailReqFields{Username: testUsername, Code: "000000"}
	req := createVerifyEmailRequest(&reqFields)
	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusBadRequest)

	// The correct code should verify the account.
	reqFields.Code = challenge.Code
	req = createVerifyEmailRequest(&reqFields)
	res = executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusOK)

	updatedMentor := models.Mentor{}
	if err := db.Where("username = ?", testUsername).First(&updatedMentor).Error; err != nil {
		t.Fatalf("Error fetching mentor: %v", err)
	}

	if !updatedMentor.EmailVerified {
		t.Error("Expected mentor to be marked as verified.")
	}
}

func TestVerifyEmailLink(t *testing.T) {
	db := setTestDB()
	defer unsetTestDB(db)

	testUsername := getTestUsername()
	student := models.Student{
		Username: testUsername,
		Name:     "Test Student",
		Email:    "student@example.com",
	}
	if err := db.Create(&student).Error; err != nil {
		t.Fatalf("Error creating test student: %v", err)
	}

	challenge, err := utils.CreateVerificationChallenge(db, testUsername, student.Email)
	if err != nil {
		t.Fatalf("Error creating verification challenge: %v", err)
	}

	req, _ := http.NewRequest(
		"GET",
		"/verify/email/?token="+challenge.Token,
		nil,
	)

	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusOK)

	if !strings.Contains(res.Body.String(), "Email verified") {
		t.Errorf("Expected a verification success page. Got `%s`.", res.Body.String())
	}

	updatedStudent := models.Student{}
	if err := db.Where("username = ?", testUsername).First(&updatedStudent).Error; err != nil {
		t.Fatalf("Error fetching student: %v", err)
	}

	if !updatedStudent.EmailVerified {
		t.Error("Expected student to be marked as verified.")
	}
}

func TestVerifyEmailLinkWithoutToken(t *testing.T) {
	db := setTestDB()
	defer unsetTestDB(db)

	req, _ := http.NewRequest("GET", "/verify/email/", nil)

	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusBadRequest)

	if !strings.Contains(res.Body.String(), "Invalid verification link") {
		t.Errorf("Expected an invalid link page. Got `%s`.", res.Body.String())
	}
}

func TestResendVerificationNoAuth(t *testing.T) {
	testRequestNoAuth(t, "POST", "/verify/email/resend/")
}

func TestVerifiedEmailRequiredMiddleware(t *testing.T) {
	db := setTestDB()
	defer unsetTestDB(db)

	t.Setenv("EMAIL_VERIFICATION_REQUIRED", "true")
	t.Setenv("EMAIL_VERIFICATION_ENABLED", "true")

	testUsername := getTestUsername()

	student := models.Student{
		Username:      testUsername,
		Email:         "student@example.com",
		EmailVerified: false,
	}
	if err := db.Create(&student).Error; err != nil {
		t.Fatalf("Error creating test student: %v", err)
	}

	app := &middleware.App{Db: db}
	handler := middleware.WithVerifiedEmail(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	execute := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest("POST", "/student/bloglink/", nil)
		ctx := context.WithValue(req.Context(), middleware.APP_CTX_KEY, app)
		ctx = context.WithValue(ctx, middleware.LOGIN_CTX_USERNAME_KEY, testUsername)
		req = req.WithContext(ctx)

		res := httptest.NewRecorder()
		handler(res, req)
		return res
	}

	res := execute()
	expectStatusCodeToBe(t, res, http.StatusForbidden)

	if err := db.Model(&models.Student{}).
		Where("id = ?", student.ID).
		Update("email_verified", true).Error; err != nil {
		t.Fatalf("Error verifying student: %v", err)
	}

	res = execute()
	expectStatusCodeToBe(t, res, http.StatusOK)
}
