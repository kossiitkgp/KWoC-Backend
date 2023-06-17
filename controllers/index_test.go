package controllers_test

import (
	"net/http"
	"testing"
)

func TestIndexController(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/", nil)
	res := executeRequest(req)

	if res.Code != http.StatusOK {
		t.Errorf("Received response code %d. Expected %d.", res.Code, http.StatusOK)
	}

	resBody := res.Body.String()
	expectedBody := "Hello from KOSS Backend!"
	if resBody != expectedBody {
		t.Errorf("Received response `%s`. Expected `%s`.", resBody, expectedBody)
	}
}
