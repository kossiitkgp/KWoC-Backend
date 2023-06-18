package controllers_test

import (
	"kwoc-backend/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Ref: https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql#h-writing-tests-for-the-api
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router := server.NewRouter(nil)

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

func testRequestNoAuth(t *testing.T, method string, path string) {
	req, _ := http.NewRequest(method, path, nil)
	res := executeRequest(req)

	expectStatusCodeToBe(t, res, http.StatusUnauthorized)
}

func testRequestInvalidAuth(t *testing.T, method string, path string) {
	req, _ := http.NewRequest(method, path, nil)
	req.Header.Add("Bearer", "Some invalid token")

	res := executeRequest(req)

	// Expect internal server error because token parsing throws an error
	expectStatusCodeToBe(t, res, http.StatusInternalServerError)
}
