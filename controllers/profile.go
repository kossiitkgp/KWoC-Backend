package controllers

import (
	"net/http"

	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"
	"gorm.io/gorm"

	"github.com/kossiitkgp/kwoc-backend/v2/models"
)

type ProfileResBodyFields struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	// `mentor` or `student`
	Type string `json:"type"`
}

// FetchProfile godoc
// @Summary		Fetches user profile
// @Description	Fetches the user's profile from the JWT, if it is valid. If invalid, returns an error.
// @Accept			plain
// @Produce		json
// @Success		200		{object}	ProfileResBodyFields	"Succesfully authenticated."
// @Failure		400		{object}	utils.HTTPMessage	"User is not registered."
// @Failure		401		{object}	utils.HTTPMessage	"JWT session token invalid."
// @Failure		500		{object}	utils.HTTPMessage	"Error parsing JWT string."
//
//	@Security		JWT
//
// @Router			/profile [get]
func FetchProfile(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	username := r.Context().Value(middleware.LOGIN_CTX_USERNAME_KEY).(utils.LoginJwtFields).Username

	// Check if the student already exists in the db
	student := models.Student{}
	tx := db.
		Table("students").
		Where("username = ?", username).
		First(&student)

	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		utils.LogErrAndRespond(r, w, tx.Error, "Database error.", http.StatusInternalServerError)
		return
	}

	student_exists := student.Username == username
	if student_exists {
		utils.RespondWithJson(r, w, ProfileResBodyFields{
			Username: student.Username,
			Name:     student.Name,
			Email:    student.Email,
			Type:     "student",
		})
		return
	}

	// Check if a mentor of the same username exists
	mentor := models.Mentor{}
	tx = db.
		Table("mentors").
		Where("username = ?", username).
		First(&mentor)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		utils.LogErrAndRespond(r, w, tx.Error, "Database error.", http.StatusInternalServerError)
		return
	}
	mentor_exists := mentor.Username == username

	if mentor_exists {
		utils.RespondWithJson(r, w, ProfileResBodyFields{
			Username: mentor.Username,
			Name:     mentor.Name,
			Email:    mentor.Email,
			Type:     "mentor",
		})
		return
	}

	utils.RespondWithHTTPMessage(r, w, http.StatusBadRequest, "User is not registered.")
}
