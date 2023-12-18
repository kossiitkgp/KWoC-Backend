package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
	"github.com/kossiitkgp/kwoc-backend/v2/models"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Mentor struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}
type Project struct {
	Id              uint     `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Tags            []string `json:"tags"`
	RepoLink        string   `json:"repo_link"`
	CommChannel     string   `json:"comm_channel"`
	ReadmeLink      string   `json:"readme_link"`
	Mentor          Mentor   `json:"mentor"`
	SecondaryMentor Mentor   `json:"secondary_mentor"`
}

func newMentor(dbMentor *models.Mentor) Mentor {
	return Mentor{
		Name:     dbMentor.Name,
		Username: dbMentor.Username,
	}
}
func newProject(dbProject *models.Project) Project {
	tags := make([]string, 0)
	if len(dbProject.Tags) != 0 {
		tags = strings.Split(dbProject.Tags, ",")
	}

	return Project{
		Id:              dbProject.ID,
		Name:            dbProject.Name,
		Description:     dbProject.Description,
		Tags:            tags,
		RepoLink:        dbProject.RepoLink,
		CommChannel:     dbProject.CommChannel,
		ReadmeLink:      dbProject.ReadmeLink,
		Mentor:          newMentor(&dbProject.Mentor),
		SecondaryMentor: newMentor(&dbProject.SecondaryMentor),
	}
}

// FetchAllProjects godoc
//
//	@Summary		Fetches all Projects
//	@Description	Fetches the public details for all the Projects
//	@Accept			plain
//	@Produce		json
//	@Success		200	{object}	[]Project	"Projects fetched successfully."
//	@Failure		500	{object}	utils.HTTPMessage	"Error fetching projects from the database."
//	@Router			/project/ [get]
func FetchAllProjects(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	var projects []models.Project

	tx := db.
		Table("projects").
		Preload("Mentor").
		Preload("SecondaryMentor").
		Where("project_status = ?", true).
		Select("id", "name", "description", "tags", "repo_link", "comm_channel", "readme_link", "mentor_id", "secondary_mentor_id").
		Find(&projects)

	if tx.Error != nil {
		utils.LogErrAndRespond(r, w, tx.Error, "Error fetching projects from the database.", http.StatusInternalServerError)
		return
	}

	var response []Project = make([]Project, 0)

	for _, project := range projects {
		response = append(response, newProject(&project))
	}

	utils.RespondWithJson(r, w, response)
}

// FetchProjectDetails godoc
//
//	@Summary		Fetches Project Details
//	@Description	Fetches the public details for the Project with the provided ID
//	@Accept			plain
//	@Produce		json
//	@Param			id	path		int			true	"Project ID"
//	@Success		200	{object}	Project	            "Project fetched successfully."
//	@Failure		400	{object}	utils.HTTPMessage	"Project id not found."
//	@Failure		400	{object}	utils.HTTPMessage	"Error parsing project id."
//	@Failure		400	{object}	utils.HTTPMessage	"Project with id `id` does not exist."
//	@Failure		400	{object}	utils.HTTPMessage	"Error: Mentor `mentor` does not own the project with ID `id`."
//	@Failure		500	{object}	utils.HTTPMessage	"Error fetching project from the database."
//	@Router			/project/{id} [get]
func FetchProjectDetails(w http.ResponseWriter, r *http.Request) {
	reqParams := mux.Vars(r)

	if reqParams["id"] == "" {
		utils.LogWarnAndRespond(r, w, "Project id not found.", http.StatusBadRequest)
		return
	}

	project_id, err := strconv.Atoi(reqParams["id"])

	if err != nil {
		utils.LogErrAndRespond(r, w, err, "Error parsing project id.", http.StatusBadRequest)
		return
	}

	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	login_username := r.Context().Value(middleware.LoginCtxKey(middleware.LOGIN_CTX_USERNAME_KEY))

	project := models.Project{}
	tx := db.
		Table("projects").
		Preload("Mentor").
		Preload("SecondaryMentor").
		Where("id = ?", project_id).
		Select("id", "name", "description", "tags", "repo_link", "comm_channel", "readme_link", "mentor_id", "secondary_mentor_id").
		First(&project)

	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		utils.LogErrAndRespond(r, w, err, "Error fetching project from the database.", http.StatusInternalServerError)
		return
	} else if tx.Error == gorm.ErrRecordNotFound {
		utils.LogWarnAndRespond(
			r,
			w,
			fmt.Sprintf("Project with id `%d` does not exist.", project_id),
			http.StatusBadRequest,
		)
		return
	}

	if project.Mentor.Username != login_username {
		utils.LogErrAndRespond(
			r,
			w,
			tx.Error,
			fmt.Sprintf("Error: Mentor `%s` does not own the project with ID `%d`.", login_username, project.ID),
			http.StatusBadRequest,
		)
		return
	}

	response := newProject(&project)
	utils.RespondWithJson(r, w, response)
}
