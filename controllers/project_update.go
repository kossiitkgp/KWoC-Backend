package controllers

import (
	"encoding/json"
	"fmt"
	"kwoc-backend/middleware"
	"kwoc-backend/models"
	"kwoc-backend/utils"
	"net/http"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UpdateProjectReqFields struct {
	// Id of the project in the database (required)
	Id uint `json:"id"`
	// Name of the project
	Name string `json:"name"`
	// Description for the project
	Description string `json:"desc"`
	// List of tags for the project
	Tags string `json:"tags"`
	// Mentor's username
	MentorUsername string `json:"username"`
	// Secondary mentor's username (if updated)
	SecondaryMentorUsername string `json:"secondaryMentor"`
	// Link to the repository of the project
	RepoLink string `json:"repoLink"`
	// Link to a communication channel/platform
	ComChannel string `json:"comChannel"`
	// Link to the project's README file
	ReadmeURL string `json:"readme"`
}

func UpdateProject(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	reqFields := UpdateProjectReqFields{}

	err := json.NewDecoder(r.Body).Decode(&reqFields)
	if err != nil {
		utils.LogErrAndRespond(r, w, err, "Error decoding request JSON body.", http.StatusBadRequest)
		return
	}

	login_username := r.Context().Value(middleware.LoginCtxKey(middleware.LOGIN_CTX_USERNAME_KEY))

	if reqFields.MentorUsername != login_username {
		log.Warn().Msgf(
			"%s %s %s\n%s %s",
			r.Method,
			r.RequestURI,
			"POSSIBLE SESSION HIJACKING.",
			fmt.Sprintf("JWT Username: %s", login_username),
			fmt.Sprintf("Given Username: %s", reqFields.MentorUsername),
		)

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Login username and mentor username do not match.")
		return
	}

	// Check if the project already exists
	project := models.Project{}
	tx := db.
		Table("projects").
		Where("id = ?", reqFields.Id).
		First(&project)

	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		utils.LogErrAndRespond(r, w, err, "Database error.", http.StatusInternalServerError)
		return
	}

	project_exists := project.ID == reqFields.Id
	if !project_exists {
		utils.LogWarnAndRespond(
			r,
			w,
			fmt.Sprintf("Error: Project `%s` does not exist.", reqFields.RepoLink),
			http.StatusBadRequest,
		)
		return
	}

	// Attempt to fetch secondary mentor from the database
	secondaryMentor := models.Mentor{}
	if reqFields.SecondaryMentorUsername != "" {
		tx = db.Table("mentors").Where("username = ?", reqFields.SecondaryMentorUsername).First(&secondaryMentor)

		if tx.Error != nil && err != gorm.ErrRecordNotFound {
			utils.LogErrAndRespond(
				r,
				w,
				err,
				fmt.Sprintf("Error fetching secondary mentor `%s`.", reqFields.SecondaryMentorUsername),
				http.StatusInternalServerError,
			)
			return
		}
	}

	updatedProj := &models.Project{
		Name:            reqFields.Name,
		Desc:            reqFields.Description,
		Tags:            reqFields.Tags,
		RepoLink:        reqFields.RepoLink,
		ComChannel:      reqFields.ComChannel,
		README:          reqFields.ReadmeURL,
		SecondaryMentor: secondaryMentor,
	}

	tx = db.
		Table("projects").
		Select("name", "desc", "tags", "repo_link", "com_channel", "readme", "secondary_mentor_id").
		Updates(updatedProj)

	if tx.Error != nil {
		utils.LogErrAndRespond(r, w, tx.Error, "Error updating the project.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Project successfully updated.")
}
