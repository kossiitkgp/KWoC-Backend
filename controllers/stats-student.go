package controllers

import (
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

type AllStudentsStats struct {
	Name     string `json:"name"`
	Username string `json:"username"`

	PrCount      uint `json:"pr_count"`
	CommitCount  uint `json:"commit_count"`
	LinesAdded   uint `json:"lines_added"`
	LinesRemoved uint `json:"lines_removed"`
}
type AllStudentsRes struct {
	Stats []AllStudentsStats
}

func AllStudents(req map[string]interface{}, r *http.Request) (interface{}, int) {
	db := utils.GetDB()
	var students []models.Student

	db.
		Table("students").
		Select("*").
		Find(&students)

	student_stats := make([]AllStudentsStats, 0)

	for _, student := range students {
		student_stats = append(
			student_stats,
			AllStudentsStats{
				Name:         student.Name,
				Username:     student.Username,
				PrCount:      student.PRCount,
				CommitCount:  student.CommitCount,
				LinesAdded:   student.AddedLines,
				LinesRemoved: student.RemovedLines,
			},
		)
	}

	response := AllStudentsRes{
		Stats: student_stats,
	}

	return response, 200
}

type OneStudentPull struct {
	Url string `json:"url"`
}
type OneStudentRepo struct {
	RepoLink string `json:"repo_link"`
	Name     string `json:"name"`
}
type OneStudentStats struct {
	Name     string `json:"name"`
	Username string `json:"username"`

	CommitCount  uint `json:"commit_count"`
	LinesAdded   uint `json:"lines_added"`
	LinesRemoved uint `json:"lines_removed"`

	Languages      []string         `json:"languages"`
	Pulls          []OneStudentPull `json:"pulls"`
	ProjectsWorked []OneStudentRepo `json:"projects_worked"`
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
		var projects_worked []OneStudentRepo = make([]OneStudentRepo, 0)

		for _, proj_id := range strings.Split(student.ProjectsWorked, ",") {
			proj := models.Project{}
			db.
				Table("projects").
				Where("id = ?", proj_id).
				First(&proj)

			projects_worked = append(
				projects_worked,
				OneStudentRepo{
					RepoLink: proj.RepoLink,
					Name:     proj.Name,
				},
			)
		}

		var pulls []OneStudentPull = make([]OneStudentPull, 0)

		for _, pull_url := range strings.Split(student.Pulls, ",") {
			pulls = append(pulls, OneStudentPull{Url: pull_url})
		}

		return OneStudentStats{
			Name:     student.Name,
			Username: student.Username,

			CommitCount:  student.CommitCount,
			LinesAdded:   student.AddedLines,
			LinesRemoved: student.RemovedLines,

			Languages:      strings.Split(student.TechWorked, ","),
			Pulls:          pulls,
			ProjectsWorked: projects_worked,
		}, 200
	} else {
		return OneStudentStats{}, 200
	}
}
