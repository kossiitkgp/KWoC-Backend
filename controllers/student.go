package controllers

import (
	"encoding/json"
	"fmt"
	"kwoc-backend/middleware"
	"net/http"

	"kwoc-backend/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type RegisterStudentReqFields struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	College  string `json:"college"`
}

func (dbHandler *DBHandler) RegisterStudent(w http.ResponseWriter, r *http.Request) {
	db := dbHandler.db
	var reqFields = RegisterStudentReqFields{}

	err := json.NewDecoder(r.Body).Decode(&reqFields)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error decoding JSON body.")

		log.Err(err).Msg("Error decoding JSON body.")
		return
	}

	// Check if the JWT login username is the same as the student's given username
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

	// Check if the student already exists in the db
	student := models.Student{}
	tx := db.
		Table("students").
		Where("username = ?", reqFields.Username).
		First(&student)

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

	student_exists := student.Username == reqFields.Username

	if student_exists {
		log.Warn().Msgf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			"Error: Student already exists.",
		)

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error: Student already exists.")
		return
	}

	// Create a db entry if the student doesn't exist
	tx = db.Create(&models.Student{
		Username: reqFields.Username,
		Name:     reqFields.Name,
		Email:    reqFields.Email,
		College:  reqFields.College,
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
