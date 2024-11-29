package server

import (
	"net/http"
	"os"

	"github.com/kossiitkgp/kwoc-backend/v2/controllers"
	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	disabled    bool
}

func getRoutes(app *middleware.App) []Route {
	return []Route{
		{
			"Index",
			"GET",
			"/api/",
			controllers.Index,
			false,
		},
		{
			"OAuth",
			"POST",
			"/oauth/",
			middleware.WrapApp(app, controllers.OAuth),
			false,
		},
		{
			"Profile",
			"GET",
			"/profile/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.FetchProfile)),
			false,
		},
		{
			"Fetch Student Details",
			"GET",
			"/student/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.GetStudentDetails)),
			false,
		},
		{
			"Fetch Mentor Details",
			"GET",
			"/mentor/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.GetMentorDetails)),
			false,
		},
		{
			"Student Registration",
			"POST",
			"/student/form/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.RegisterStudent)),
			os.Getenv("REGISTRATIONS_OPEN") == "true",
		},
		{
			"Update Student Details",
			"PUT",
			"/student/form/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.UpdateStudentDetails)),
			false,
		},
		{
			"Student Blog Submission",
			"POST",
			"/student/bloglink/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.StudentBlogLink)),
			os.Getenv("REPORT_SUBMISSION_OPEN") == "true",
		},
		{
			"Student Dashboard",
			"GET",
			"/student/dashboard/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.FetchStudentDashboard)),
			false,
		},
		{
			"Mentor Registration",
			"POST",
			"/mentor/form/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.RegisterMentor)),
			os.Getenv("REGISTRATIONS_OPEN") == "true",
		},
		{
			"Update Mentor Details",
			"PUT",
			"/mentor/form/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.UpdateMentorDetails)),
			false,
		},
		{
			"Mentor Dashboard",
			"GET",
			"/mentor/dashboard/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.FetchMentorDashboard)),
			false,
		},
		{
			"HealthCheck",
			"GET",
			"/healthcheck/",
			middleware.WrapApp(app, controllers.HealthCheck),
			false,
		},
		{
			"Ping",
			"GET",
			"/healthcheck/ping/",
			controllers.Ping,
			false,
		},
		{
			"Project Registration",
			"POST",
			"/project/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.RegisterProject)),
			os.Getenv("REGISTRATIONS_OPEN") == "true",
		},
		{
			"Fetch All Projects",
			"GET",
			"/project/",
			middleware.WrapApp(app, controllers.FetchAllProjects),
			false,
		}, {
			"Get All Projects For Admins",
			"GET",
			"/project/all",
			middleware.WithLogin(middleware.WrapApp(app, controllers.OrgFetchAllProjectDetails)),
			false,
		},
		{
			"Update Project Details",
			"PUT",
			"/project/",
			middleware.WithLogin(middleware.WrapApp(app, controllers.UpdateProject)),
			false,
		},
		{
			"Fetch Project Details",
			"GET",
			"/project/{id}",
			middleware.WithLogin(middleware.WrapApp(app, controllers.FetchProjectDetails)),
			false,
		},
		{
			"Fetch All Students Stats",
			"GET",
			"/stats/students/",
			middleware.WrapApp(app, controllers.FetchAllStudentStats),
			false,
		},
		{
			"Fetch Overall Stats",
			"GET",
			"/stats/overall/",
			middleware.WrapApp(app, controllers.FetchOverallStats),
			false,
		},
		{
			"Fetch Project Stats",
			"GET",
			"/stats/projects/",
			middleware.WrapApp(app, controllers.FetchAllProjectStats),
			false,
		},
	}
}
