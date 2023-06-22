package controllers

import (
	"kwoc-backend/middleware"
	"kwoc-backend/models"
	"kwoc-backend/utils"
	"net/http"
)

type FetchAllProjMentor struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}
type FetchAllProjProject struct {
	Name            string             `json:"name"`
	Desc            string             `json:"desc"`
	Tags            string             `json:"tags"`
	RepoLink        string             `json:"repo_link"`
	ComChannel      string             `json:"com_channel"`
	Mentor          FetchAllProjMentor `json:"mentor"`
	SecondaryMentor FetchAllProjMentor `json:"secondary_mentor"`
}

type FetchAllProjRes []FetchAllProjProject

func newFetchAllProjMentor(mentor *models.Mentor) FetchAllProjMentor {
	return FetchAllProjMentor{
		Name:     mentor.Name,
		Username: mentor.Username,
	}
}
func newFetchAllProjProject(project *models.Project) FetchAllProjProject {
	return FetchAllProjProject{
		Name:            project.Name,
		Desc:            project.Desc,
		Tags:            project.Tags,
		RepoLink:        project.RepoLink,
		ComChannel:      project.ComChannel,
		Mentor:          newFetchAllProjMentor(&project.Mentor),
		SecondaryMentor: newFetchAllProjMentor(&project.SecondaryMentor),
	}
}

func FetchAllProjects(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	var projects []models.Project

	tx := db.
		Table("projects").
		Preload("Mentor").
		Preload("SecondaryMentor").
		Where("project_status = ?", true).
		Select("name", "desc", "tags", "repo_link", "com_channel", "mentor_id", "secondary_mentor_id").
		Find(&projects)

	if tx.Error != nil {
		utils.LogErrAndRespond(r, w, tx.Error, "Error fetching projects from the database.", http.StatusInternalServerError)
		return
	}

	var response FetchAllProjRes = make(FetchAllProjRes, 0)

	for _, project := range projects {
		response = append(response, newFetchAllProjProject(&project))
	}

	utils.RespondWithJson(r, w, response)
}
