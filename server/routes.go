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

func getRoutes(dbHandler *controllers.DBHandler) []Route {
	return []Route{
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
			dbHandler.OAuth,
		},
		{
			"Student Registration",
			"POST",
			"/student/form/",
			middleware.WithLogin(dbHandler.RegisterStudent),
		},
		{
			"Mentor Registration",
			"POST",
			"/mentor/form/",
			middleware.WithLogin(dbHandler.RegisterMentor),
		},
	}
}
