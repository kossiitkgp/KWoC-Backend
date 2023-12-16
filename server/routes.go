package server

import (
	"net/http"

	"github.com/kossiitkgp/kwoc-backend/v2/controllers"
	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
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
			"Profile",
			"GET",
			"/profile/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.FetchProfile)),
		},
		{
			"Fetch Student Details",
			"GET",
			"/student/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.GetStudentDetails)),
		},
		{
			"Fetch Mentor Details",
			"GET",
			"/mentor/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.GetMentorDetails)),
		},
		// {
		// 	"Student Registration",
		// 	"POST",
		// 	"/student/form/",
		// 	middleware.WithLogin(middleware.WrapApp(app, controllers.RegisterStudent)),
		// },
		{
			"Update Student Details",
			"PUT",
			"/student/form/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.UpdateStudentDetails)),
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
		// {
		// 	"Mentor Registration",
		// 	"POST",
		// 	"/mentor/form/",
		// 	middleware.WithLogin(middleware.WrapApp(app, controllers.RegisterMentor)),
		// },
		{
			"Update Mentor Details",
			"PUT",
			"/mentor/form/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.UpdateMentorDetails)),
		},
		{
			"Fetch All Mentors",
			"GET",
			"/mentor/all/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.FetchAllMentors)),
		},
		{
			"Mentor Dashboard",
			"GET",
			"/mentor/dashboard/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.FetchMentorDashboard)),
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
		// {
		// 	"Project Registration",
		// 	"POST",
		// 	"/project/",
		// 	middleware.WithLogin(middleware.WrapApp(app, controllers.RegisterProject)),
		// },
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
			middleware.WithLogin(middleware.WrapApp(app, controllers.FetchProjectDetails)),
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
