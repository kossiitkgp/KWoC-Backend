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

type RegisterProjectReqFields struct {
	// Name of the project
	Name string `json:"name"`
	// Description for the project
	Description string `json:"desc"`
	// List of tags for the project
	Tags string `json:"tags"`
	// Mentor's username
	MentorUsername string `json:"username"`
	// Secondary mentor's username
	SecondaryMentorUsername string `json:"secondaryMentor"`
	// Link to the repository of the project
	RepoLink string `json:"repoLink"`
	// Link to a communication channel/platform
	ComChannel string `json:"comChannel"`
	// Link to the project's README file
	ReadmeURL string `json:"readme"`
}

func RegisterProject(w http.ResponseWriter, r *http.Request) {
	reqFields := RegisterProjectReqFields{}

	err := json.NewDecoder(r.Body).Decode(&reqFields)

	if err != nil {
		log.Err(err).Msgf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			"Error decoding request JSON body.",
		)

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error decoding request JSON body.")
		return
	}

	login_username := r.Context().Value(middleware.LOGIN_CTX_USERNAME_KEY)

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

	db, err := utils.GetDB()
	if err != nil {
		log.Err(err).Msg("Error connecting to the database.")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error connecting to the database.")
		return
	}

	// Fetch primary mentor from the database
	mentor := models.Mentor{}
	tx := db.Table("mentors").Where("username = ?", reqFields.MentorUsername).First(&mentor)

	if tx.Error != nil {
		log.Err(err).Msgf("Error fetching mentor `%s`.", reqFields.MentorUsername)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error fetching mentor.")
		return
	}

	// Attempt to fetch secondary mentor from the database
	secondaryMentor := models.Mentor{}
	if reqFields.SecondaryMentorUsername != "" {
		tx = db.Table("mentors").Where("username = ?", reqFields.SecondaryMentorUsername).First(&secondaryMentor)

		if tx.Error != nil && err != gorm.ErrRecordNotFound {
			log.Err(err).Msgf("Error fetching secondary mentor `%s`.", reqFields.SecondaryMentorUsername)

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error fetching secondary mentor.")
			return
		}
	}

	tx = db.Create(models.Project{
		Name:            reqFields.Name,
		Desc:            reqFields.Description,
		Tags:            reqFields.Tags,
		RepoLink:        reqFields.RepoLink,
		ComChannel:      reqFields.ComChannel,
		README:          reqFields.ReadmeURL,
		Mentor:          mentor,
		SecondaryMentor: secondaryMentor,
	})

	if tx.Error != nil {
		log.Err(err).Msgf("Error creating project in the database.")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error registering project.")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Success.")
}
