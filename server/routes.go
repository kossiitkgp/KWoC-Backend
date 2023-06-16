package server

import (
	"kwoc-backend/controllers"
	"kwoc-backend/middleware"
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
	{
		"Student Registration",
		"POST",
		"/student/form/",
		middleware.WithLogin(controllers.RegisterStudent),
	},
	{
		"Mentor Registration",
		"POST",
		"/mentor/form/",
		middleware.WithLogin(controllers.RegisterMentor),
	},
}
