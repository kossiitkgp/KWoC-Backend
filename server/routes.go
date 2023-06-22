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

func getRoutes(app *middleware.App) []Route {
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
			middleware.WrapApp(app, controllers.OAuth),
		},
		{
			"Student Registration",
			"POST",
			"/student/form/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.RegisterStudent)),
		},
		{
			"Mentor Registration",
			"POST",
			"/mentor/form/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.RegisterMentor)),
		},
		{
			"HealthCheck",
			"GET",
			"/healthcheck/",
			middleware.WrapApp(app, controllers.HealthCheck),
		},
		{
			"Ping",
			"GET",
			"/healthcheck/ping/",
			controllers.Ping,
		},
		{
			"Project Registration",
			"POST",
			"/project/add/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.RegisterProject)),
		},
	}
}
