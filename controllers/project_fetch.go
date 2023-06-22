package controllers

import (
	"fmt"
	"kwoc-backend/middleware"
	"kwoc-backend/models"
	"kwoc-backend/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type FetchProjMentor struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}
type FetchProjProject struct {
	Name            string          `json:"name"`
	Desc            string          `json:"desc"`
	Tags            string          `json:"tags"`
	RepoLink        string          `json:"repo_link"`
	ComChannel      string          `json:"com_channel"`
	Mentor          FetchProjMentor `json:"mentor"`
	SecondaryMentor FetchProjMentor `json:"secondary_mentor"`
}

type FetchAllProjRes []FetchProjProject

func newFetchProjMentor(mentor *models.Mentor) FetchProjMentor {
	return FetchProjMentor{
		Name:     mentor.Name,
		Username: mentor.Username,
	}
}
func newFetchProjProject(project *models.Project) FetchProjProject {
	return FetchProjProject{
		Name:            project.Name,
		Desc:            project.Desc,
		Tags:            project.Tags,
		RepoLink:        project.RepoLink,
		ComChannel:      project.ComChannel,
		Mentor:          newFetchProjMentor(&project.Mentor),
		SecondaryMentor: newFetchProjMentor(&project.SecondaryMentor),
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
		response = append(response, newFetchProjProject(&project))
	}

	utils.RespondWithJson(r, w, response)
}

func FetchProjDetails(w http.ResponseWriter, r *http.Request) {
	reqParams := mux.Vars(r)

	if reqParams["id"] == "" {
		utils.LogWarnAndRespond(r, w, "Project id not found.", http.StatusBadRequest)
		return
	}

	proj_id, err := strconv.Atoi(reqParams["id"])

	if err != nil {
		utils.LogErrAndRespond(r, w, err, "Error parsing project id.", http.StatusBadRequest)
		return
	}

	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	project := models.Project{}
	tx := db.
		Table("projects").
		Preload("Mentor").
		Preload("SecondaryMentor").
		Where("project_status = ?", true).
		Where("id = ?", proj_id).
		Select("name", "desc", "tags", "repo_link", "com_channel", "mentor_id", "secondary_mentor_id").
		First(&project)

	if tx.Error != nil {
		utils.LogErrAndRespond(r, w, err, "Error fetching project from the database.", http.StatusInternalServerError)
	}

	if int(project.ID) != proj_id {
		utils.LogWarnAndRespond(
			r,
			w,
			fmt.Sprintf("Project with id `%d` does not exist.", proj_id),
			http.StatusBadRequest,
		)
		return
	}

	response := newFetchProjProject(&project)
	utils.RespondWithJson(r, w, response)
}
