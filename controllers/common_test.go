package controllers_test

import (
	"kwoc-backend/server"
	"net/http"
	"net/http/httptest"
)

// Ref: https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql#h-writing-tests-for-the-api
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router := server.NewRouter()

	router.ServeHTTP(rr, req)

	return rr
}
