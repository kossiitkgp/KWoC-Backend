package server

import (
	"kwoc-backend/controllers"
	"kwoc-backend/utils"
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
		utils.WithLogin(controllers.Index),
	},
}
