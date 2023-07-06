package controllers

import (
	"fmt"
	"kwoc-backend/middleware"
	"kwoc-backend/models"
	"kwoc-backend/utils"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type RegisterStudentReqFields struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	College  string `json:"college"`
}

type StudentBlogLinkReqFields struct {
	Username string `json:"username"`
	BlogLink string `json:"blog_link"`
}

type ProjectDashboard struct {
	Name     string `json:"name"`
	RepoLink string `json:"repo_link"`
}

type StudentDashboard struct {
	Name           string `json:"name"`
	Username       string `json:"username"`
	College        string `json:"college"`
	PassedMidEvals bool   `json:"passed_mid_evals"`
	PassedEndEvals bool   `json:"passed_end_evals"`

	CommitCount  uint `json:"commit_count"`
	PullCount    uint `json:"pull_count"`
	LinesAdded   uint `json:"lines_added"`
	LinesRemoved uint `json:"lines_removed"`

	LanguagesUsed  []string           `json:"languages_used"`
	ProjectsWorked []ProjectDashboard `json:"projects_worked"`
}

func RegisterStudent(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db
	reqFields := RegisterStudentReqFields{}

	err := utils.DecodeJSONBody(r, &reqFields)
	if err != nil {
		utils.LogErrAndRespond(r, w, err, "Error decoding JSON body.", http.StatusBadRequest)
		return
	}

	// Check if the JWT login username is the same as the student's given username
	login_username := r.Context().Value(middleware.LOGIN_CTX_USERNAME_KEY).(string)

	if reqFields.Username != login_username {
		utils.LogWarn(
			r,
			fmt.Sprintf(
				"POSSIBLE SESSION HIJACKING\nJWT Username: %s, Given Username: %s",
				login_username,
				reqFields.Username,
			),
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
		utils.LogErrAndRespond(r, w, err, "Database error.", http.StatusInternalServerError)
		return
	}

	student_exists := student.Username == reqFields.Username

	if student_exists {
		utils.LogWarnAndRespond(
			r,
			w,
			fmt.Sprintf("Student `%s` already exists.", student.Username),
			http.StatusBadRequest,
		)
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
		utils.LogErrAndRespond(r, w, err, "Database error.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Student registration successful.")
}

func StudentBlogLink(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db
	reqFields := StudentBlogLinkReqFields{}

	err := utils.DecodeJSONBody(r, &reqFields)
	if err != nil {
		utils.LogErrAndRespond(r, w, err, "Error decoding JSON body.", http.StatusBadRequest)
		return
	}

	// Check if the JWT login username is the same as the student's given username
	login_username := r.Context().Value(middleware.LOGIN_CTX_USERNAME_KEY).(string)

	if reqFields.Username != login_username {
		utils.LogWarn(
			r,
			fmt.Sprintf(
				"POSSIBLE SESSION HIJACKING\nJWT Username: %s, Given Username: %s",
				login_username,
				reqFields.Username,
			),
		)

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Login username and given username do not match.")
		return
	}

	// Check if the student exists in the db
	student := models.Student{}
	tx := db.
		Table("students").
		Where("username = ?", reqFields.Username).
		First(&student)

	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		utils.LogErrAndRespond(r, w, err, "Database error.", http.StatusInternalServerError)
		return
	}

	if tx.Error == gorm.ErrRecordNotFound {
		utils.LogWarnAndRespond(
			r,
			w,
			fmt.Sprintf("Student `%s` does not exists.", reqFields.Username),
			http.StatusBadRequest,
		)
		return
	}
	tx = tx.Update("BlogLink", reqFields.BlogLink)

	if tx.Error != nil {
		utils.LogErrAndRespond(
			r,
			w,
			tx.Error,
			fmt.Sprintf("Error updating BlogLink for `%s`.", student.Username),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "BlogLink successfully updated.")
}

func CreateStudentDashboard(modelStudent models.Student, db *gorm.DB) StudentDashboard {
	var projects []ProjectDashboard
	for _, proj_id := range strings.Split(modelStudent.ProjectsWorked, ",") {
		var project ProjectDashboard
		db.Table("projects").
			Where("id = ?", proj_id).
			Select("name", "repo_link").
			First(&project)
		projects = append(projects, project)
	}
	languages := strings.Split(modelStudent.LanguagesUsed, ",")
	return StudentDashboard{
		Name:           modelStudent.Name,
		Username:       modelStudent.Username,
		College:        modelStudent.College,
		PassedMidEvals: modelStudent.PassedMidEvals,
		PassedEndEvals: modelStudent.PassedEndEvals,
		CommitCount:    modelStudent.CommitCount,
		PullCount:      modelStudent.PullCount,
		LinesAdded:     modelStudent.LinesAdded,
		LinesRemoved:   modelStudent.LinesRemoved,
		LanguagesUsed:  languages,
		ProjectsWorked: projects,
	}
}

func FetchStudentDashboard(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	var modelStudent models.Student

	login_username := r.Context().Value(middleware.LoginCtxKey(middleware.LOGIN_CTX_USERNAME_KEY))
	tx := db.
		Table("students").
		Where("username = ?", login_username).
		First(&modelStudent)

	if tx.Error == gorm.ErrRecordNotFound {
		utils.LogErrAndRespond(
			r,
			w,
			tx.Error,
			fmt.Sprintf("Student `%s` does not exists.", login_username),
			http.StatusBadRequest,
		)
		return
	}
	if tx.Error != nil {
		utils.LogErrAndRespond(
			r,
			w,
			tx.Error,
			fmt.Sprintf("Database Error fetching student with username `%s`", login_username),
			http.StatusInternalServerError,
		)
		return
	}
	student := CreateStudentDashboard(modelStudent, db)
	utils.RespondWithJson(r, w, student)
}
