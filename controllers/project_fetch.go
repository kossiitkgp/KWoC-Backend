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
	return Project{
		Id:              dbProject.ID,
		Name:            dbProject.Name,
		Description:     dbProject.Description,
		Tags:            strings.Split(dbProject.Tags, ","),
		RepoLink:        dbProject.RepoLink,
		CommChannel:     dbProject.CommChannel,
		ReadmeLink:      dbProject.ReadmeLink,
		Mentor:          newMentor(&dbProject.Mentor),
		SecondaryMentor: newMentor(&dbProject.SecondaryMentor),
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

	project := models.Project{}
	tx := db.
		Table("projects").
		Preload("Mentor").
		Preload("SecondaryMentor").
		Where("project_status = ?", true).
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

	response := newProject(&project)
	utils.RespondWithJson(r, w, response)
}
