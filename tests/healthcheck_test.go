package controllers_test

import (
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	req, _ := http.NewRequest("GET", "/healthcheck/ping/", nil)
	res := executeRequest(req, nil)

	expectStatusCodeToBe(t, res, http.StatusOK)
	expectResponseBodyToBe(t, res, "pong")
}

func TestHealthCheck(t *testing.T) {
	db := setTestDB()
	defer unsetTestDB(db)

	req, _ := http.NewRequest("GET", "/healthcheck/", nil)
	res := executeRequest(req, db)

	expectStatusCodeToBe(t, res, http.StatusOK)
	expectResponseBodyToBe(t, res, "OK")
}
