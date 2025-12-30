package controllers_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kossiitkgp/kwoc-backend/v2/models"
	"github.com/kossiitkgp/kwoc-backend/v2/server"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"

	"gorm.io/gorm"
)

// Ref: https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql#h-writing-tests-for-the-api
func executeRequest(req *http.Request, db *gorm.DB) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router := server.NewRouter(db, true)

	router.ServeHTTP(rr, req)

	return rr
}

func expectStatusCodeToBe(t *testing.T, res *httptest.ResponseRecorder, expectedCode int) {
	if res.Code != expectedCode {
		t.Errorf("Expected status code %d. Got %d.", expectedCode, res.Code)
	}
}

func expectResponseBodyToBe(t *testing.T, res *httptest.ResponseRecorder, expectedBody string) {
	resBody := res.Body.String()

	if resBody != expectedBody {
		t.Errorf("Expected response `%s`. Got `%s`.", expectedBody, resBody)
	}
}

func expectResponseJSONBodyToBe[T comparable](t *testing.T, res *httptest.ResponseRecorder, expectedBody T) {
	var resBody T
	err := json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		t.Errorf("Error decoding `%s`.", res.Body.String())
	}

	if resBody != expectedBody {
		t.Errorf("Expected response `%v`. Got `%v`.", expectedBody, resBody)
	}
}

func setTestDB() *gorm.DB {
	_ = godotenv.Load("../.env")
	db, _ := utils.GetDB()
	_ = utils.MigrateModels(db)

	return db
}

func unsetTestDB(db *gorm.DB) {
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.Project{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.Student{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.Stats{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.Mentor{})
}

func setTestJwtSecretKey() {
	os.Setenv("JWT_SECRET_KEY", fmt.Sprintf("testkey%d", rand.Int()))
}

func unsetTestJwtSecretKey() {
	os.Unsetenv("JWT_SECRET_KEY")
}

func getTestUsername() string {
	return fmt.Sprintf("testuser%d", rand.Int())
}

func testRequestNoAuth(t *testing.T, method string, path string) {
	req, _ := http.NewRequest(method, path, nil)
	res := executeRequest(req, nil)

	expectStatusCodeToBe(t, res, http.StatusUnauthorized)
}

func testRequestInvalidAuth(t *testing.T, method string, path string) {
	req, _ := http.NewRequest(method, path, nil)
	req.Header.Add("Bearer", "Some invalid token")

	res := executeRequest(req, nil)

	// Expect internal server error because token parsing throws an error
	expectStatusCodeToBe(t, res, http.StatusInternalServerError)
}
