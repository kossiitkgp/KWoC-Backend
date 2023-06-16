package server

import (
	"kwoc-backend/controllers"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes []Route = []Route{
	{
		"Index",
		"GET",
		"/api/",
		controllers.Index,
	},
	{
		"OAuth",
		"POST",
		"/oauth/",
		controllers.OAuth,
	},
}
