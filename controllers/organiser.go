package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
	"github.com/kossiitkgp/kwoc-backend/v2/models"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"
)

type AcceptRejectProject struct {
	// Id of the project in the database (required)
	Id uint `json:"id"`
	// Status to be set of the project
	ProjectStatus bool `json:"project_status"`
	// Status Remark to be set of the project
	StatusRemark string `json:"status_remark"`
}

type ProjectOrg struct {
	ProjectStatus bool   `json:"project_status"`
	StatusRemark  string `json:"status_remark"`
	Project
}

func newProjectOrg(dbProject *models.Project) ProjectOrg {
	tags := make([]string, 0)
	if len(dbProject.Tags) != 0 {
		tags = strings.Split(dbProject.Tags, ",")
	}

	return ProjectOrg{
		Project: Project{
			Id:              dbProject.ID,
			Name:            dbProject.Name,
			Description:     dbProject.Description,
			Tags:            tags,
			RepoLink:        dbProject.RepoLink,
			CommChannel:     dbProject.CommChannel,
			ReadmeLink:      dbProject.ReadmeLink,
			Mentor:          newMentor(&dbProject.Mentor),
			SecondaryMentor: newMentor(&dbProject.SecondaryMentor),
		},
		ProjectStatus: dbProject.ProjectStatus,
		StatusRemark:  dbProject.StatusRemark,
	}
}

func OrgFetchAllProjectDetails(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db
	user_details := r.Context().Value(middleware.LOGIN_CTX_USERNAME_KEY).(utils.LoginJwtFields)

	if user_details.UserType != "organiser" {
		utils.LogErrAndRespond(r, w, nil, fmt.Sprintf("Error '%s' is not an organiser", user_details.Username), 400)
		return
	}

	var projects []models.Project

	tx := db.
		Table("projects").
		Preload("Mentor").
		Preload("SecondaryMentor").
		Select("id", "name", "description", "tags", "repo_link", "comm_channel", "readme_link", "mentor_id", "secondary_mentor_id", "project_status", "status_remark", "pull_count").
		Find(&projects)

	if tx.Error != nil {
		utils.LogErrAndRespond(r, w, tx.Error, "Error fetching projects from the database.", http.StatusInternalServerError)
		return
	}

	var response []ProjectOrg = make([]ProjectOrg, 0)

	for _, project := range projects {
		response = append(response, newProjectOrg(&project))
	}

	utils.RespondWithJson(r, w, response)
}
