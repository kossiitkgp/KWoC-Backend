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
			"Student Blog Submission",
			"POST",
			"/student/bloglink/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.StudentBlogLink)),
		},
		{
			"Student Dashboard",
			"GET",
			"/student/dashboard/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.FetchStudentDashboard)),
		},
		{
			"Mentor Registration",
			"POST",
			"/mentor/form/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.RegisterMentor)),
		},
		{
			"Fetch All Mentors",
			"GET",
			"/mentor/all/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.FetchAllMentors)),
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
			"/project/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.RegisterProject)),
		},
		{
			"Fetch All Projects",
			"GET",
			"/project/",
			middleware.WrapApp(app, controllers.FetchAllProjects),
		},
		{
			"Update Project Details",
			"PUT",
			"/project/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.UpdateProject)),
		},
		{
			"Fetch Project Details",
			"GET",
			"/project/{id}",
			middleware.WrapApp(app, controllers.FetchProjectDetails),
		},
		{
			"Fetch All Students Stats",
			"GET",
			"/stats/students/",
			middleware.WrapApp(app, controllers.FetchAllStudentStats),
		},
		{
			"Fetch Overall Stats",
			"GET",
			"/stats/overall/",
			middleware.WrapApp(app, controllers.FetchOverallStats),
		},
		{
			"Fetch Project Stats",
			"GET",
			"/stats/projects/",
			middleware.WrapApp(app, controllers.FetchAllProjectStats),
		},
	}
}
