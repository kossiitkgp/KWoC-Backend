package controllers

import (
	"encoding/json"
	"fmt"
	"kwoc-backend/middleware"
	"kwoc-backend/utils"
	"net/http"

	"github.com/kossiitkgp/kwoc-db-models/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type RegisterMentorReqFields struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

func RegisterMentor(w http.ResponseWriter, r *http.Request) {
	var reqFields = RegisterMentorReqFields{}

	err := json.NewDecoder(r.Body).Decode(&reqFields)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error decoding JSON body.")

		log.Err(err).Msg("Error decoding JSON body.")
		return
	}

	// Check if the JWT login username is the same as the mentor's given username
	login_username := r.Context().Value(middleware.LoginCtxKey(middleware.LOGIN_CTX_USERNAME_KEY)).(string)

	if reqFields.Username != login_username {
		log.Warn().Msgf(
			"%s %s %s\n%s %s",
			r.Method,
			r.RequestURI,
			"POSSIBLE SESSION HIJACKING.",
			fmt.Sprintf("JWT Username: %s", login_username),
			fmt.Sprintf("Given Username: %s", reqFields.Username),
		)

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Login username and given username do not match.")
		return
	}

	db, err := utils.GetDB()
	if err != nil {
		log.Err(err).Msgf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			"Error connecting to database.",
		)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error connecting to database.")
		return
	}

	// Check if the mentor already exists in the db
	mentor := models.Mentor{}
	tx := db.
		Table("mentors").
		Where("username = ?", reqFields.Username).
		First(&mentor)

	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		log.Err(err).Msgf(
			"%s %s %s %v",
			r.Method,
			r.RequestURI,
			"Database error.",
			tx.Error,
		)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Database error.")
		return
	}

	mentor_exists := mentor.Username == reqFields.Username

	if mentor_exists {
		log.Warn().Msgf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			"Error: Mentor already exists.",
		)

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error: Mentor already exists.")
		return
	}

	// Create a db entry if the mentor doesn't exist
	tx = db.Create(&models.Mentor{
		Username: reqFields.Username,
		Name:     reqFields.Name,
		Email:    reqFields.Email,
	})

	if tx.Error != nil {
		log.Err(err).Msgf(
			"%s %s %s %v",
			r.Method,
			r.RequestURI,
			"Database error.",
			tx.Error,
		)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Database error.")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Success.")
}
