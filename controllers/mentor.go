package controllers

import (
	"net/http"

	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
	"github.com/kossiitkgp/kwoc-backend/v2/models"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"
	"gorm.io/gorm"
)

/*
PUBLIC STRUCT
Only username + display name
*/
type PublicMentor struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

/*
SINGLE SOURCE OF TRUTH
This must exist ONLY ONCE in the entire codebase
*/
func NewPublicMentor(m *models.Mentor) PublicMentor {
	if m == nil {
		return PublicMentor{}
	}

	return PublicMentor{
		Username:    m.Username,
		DisplayName: m.Name,
	}
}

/*
Existing mentor handlers below
(NO LOGIC CHANGED)
*/

func FetchMentors(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	var mentors []models.Mentor
	tx := db.Table("mentors").Find(&mentors)

	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		utils.LogErrAndRespond(r, w, tx.Error, "Database error.", http.StatusInternalServerError)
		return
	}

	publicMentors := make([]PublicMentor, 0)
	for i := range mentors {
		publicMentors = append(publicMentors, NewPublicMentor(&mentors[i]))
	}

	utils.RespondWithJson(r, w, publicMentors)
}
