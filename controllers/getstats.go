package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"kwoc20-backend/models"

	"kwoc20-backend/utils"

	"github.com/gorilla/mux"
)

func CheckStudent(req map[string]interface{}, r *http.Request) (interface{}, int) {
	db := utils.GetDB()
	params := mux.Vars(r)

	student := models.Student{}
	db.
		Table("students").
		Where("username = ?", params["username"]).
		First(&student)

	student_exists := student.Username == params["username"]

	if student_exists {
		return "true", 200
	} else {
		return "false", 200
	}
}

type AllStudentStat struct {
	Name     string
	Username string
	Prs      string
	Commits  uint
	Lines    string
}
type AllStudentsRes struct {
	Stats []AllStudentStat
}

func AllStudents(req map[string]interface{}, r *http.Request) (interface{}, int) {
	db := utils.GetDB()
	var students []models.Student

	db.
		Table("students").
		Select("*").
		Find(&students)

	student_stats := make([]AllStudentStat, 0)

	for _, student := range students {
		student_stats = append(
			student_stats,
			AllStudentStat{
				Name:     student.Name,
				Username: student.Username,
				Prs:      fmt.Sprintf("%d", student.PRCount),
				Commits:  student.CommitCount,
				Lines:    fmt.Sprintf("+%d/-%d", student.AddedLines, student.RemovedLines),
			},
		)
	}

	response := AllStudentsRes{
		Stats: student_stats,
	}

	return response, 200
}

type OneStudentCommit struct {
	Html_url string
	Message  string
}
type OneStudentPull struct {
	Html_url     string
	Title        string
	RepoOwner    string
	LinesAdded   string
	RepoName     string
	LinesRemoved string
}
type OneStudentStat struct {
	Name         string
	Username     string
	CommitCount  uint
	Languages    []string
	Pulls        []OneStudentPull
	LinesAdded   uint
	LinesRemoved uint
}

func OneStudent(req map[string]interface{}, r *http.Request) (interface{}, int) {
	db := utils.GetDB()
	params := mux.Vars(r)

	student := models.Student{}
	db.
		Table("students").
		Where("username = ?", params["username"]).
		First(&student)

	student_exists := student.Username == params["username"]

	if student_exists {
		return OneStudentStat{
			Name:         student.Name,
			Username:     student.Username,
			CommitCount:  student.CommitCount,
			Languages:    strings.Split(student.TechWorked, ","),
			Pulls:        make([]OneStudentPull, 0),
			LinesAdded:   student.AddedLines,
			LinesRemoved: student.RemovedLines,
		}, 200
	} else {
		return OneStudentStat{}, 200
	}
}

type Project struct {
	ID       uint
	RepoLink string
	Branch   string

	LastPullDate string

	CommitCount  uint
	PRCount      uint
	AddedLines   uint
	RemovedLines uint
}

func GetAllProjects(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	w.WriteHeader(200)
	var projects []Project
	db.
		Table("projects").
		Where("project_status = ?", "true").
		Select(
			"id", "repo_link", "branch",
			"last_pull_date", "commit_count",
			"pr_count", "added_lines", "removed_lines",
		).
		Find(&projects)
	str := fmt.Sprintf("%+v", projects)
	_, err := w.Write([]byte(str))

	if err != nil {
		utils.LogErr(r, err, "Write failed.")
	}
}

func OneMentor(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	w.WriteHeader(200)
	params := mux.Vars(r)
	var mentor []Project
	db.
		Table("projects").
		Where("project_status = ?", "true").
		Select(
			"id", "repo_link", "branch",
			"last_pull_date", "commit_count",
			"pr_count", "added_lines", "removed_lines",
		).
		Find(&mentor, params["Mentor.Username"])
	str := fmt.Sprintf("%+v", mentor)
	_, err := w.Write([]byte(str))

	if err != nil {
		utils.LogErr(r, err, "Write failed.")
	}
}
