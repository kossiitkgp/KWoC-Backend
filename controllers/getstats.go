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
	Name     string
	Username string

	PrCount      uint
	CommitCount  uint
	LinesAdded   uint
	LinesRemoved uint
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
	Url string
}
type OneStudentRepo struct {
	RepoLink string
	Name     string
}
type OneStudentStats struct {
	Name     string
	Username string

	CommitCount  uint
	LinesAdded   uint
	LinesRemoved uint

	Languages      []string
	Pulls          []OneStudentPull
	ProjectsWorked []OneStudentRepo
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

type AllProjectsProject struct {
	Name string
	Link string

	CommitCount  uint
	PrCount      uint
	LinesAdded   uint
	LinesRemoved uint

	Contributors []string
}
type AllProjectsRes struct {
	Stats []AllProjectsProject
}

func GetAllProjects(req map[string]interface{}, r *http.Request) (interface{}, int) {
	db := utils.GetDB()
	var projects []models.Project

	db.
		Table("projects").
		Where("project_status = ?", "1").
		Select("*").
		Find(&projects)

	response := make([]AllProjectsProject, 0)

	for _, project := range projects {
		response = append(
			response,
			AllProjectsProject{
				Name: project.Name,
				Link: project.RepoLink,

				CommitCount:  project.CommitCount,
				LinesAdded:   project.AddedLines,
				LinesRemoved: project.RemovedLines,

				Contributors: strings.Split(project.Contributors, ","),
			},
		)
	}

	return AllProjectsRes{Stats: response}, 200
}

type OneMentorProjPull struct {
	Url string
}
type OneMentorProj struct {
	Name     string
	RepoLink string

	CommitCount  uint
	LinesAdded   uint
	LinesRemoved uint

	Contributors []string // Array of usernames of students who contributed
	Pulls        []OneMentorProjPull
}
type OneMentorRes struct {
	Projects []OneMentorProj
}

func OneMentor(req map[string]interface{}, r *http.Request) (interface{}, int) {
	db := utils.GetDB()
	params := mux.Vars(r)
	username := params["Mentor.Username"]

	var mentor models.Mentor

	db.
		Table("mentors").
		Where("username = ?", username).
		Select("*").
		First(&mentor)

	if mentor.Username == username {
		mentor_id := mentor.ID
		var projects []models.Project

		db.
			Table("projects").
			Where("mentor_id = ? OR secondary_mentor_id = ?", mentor_id, mentor_id).
			Find(&projects)

		mentor_stats := make([]OneMentorProj, 0)

		for _, project := range projects {
			var proj_pulls []OneMentorProjPull = make([]OneMentorProjPull, 0)

			for _, pull_url := range strings.Split(project.Pulls, ",") {
				proj_pulls = append(proj_pulls, OneMentorProjPull{Url: pull_url})
			}

			mentor_stats = append(
				mentor_stats,
				OneMentorProj{
					Name:     project.Name,
					RepoLink: project.RepoLink,

					CommitCount:  project.CommitCount,
					LinesAdded:   project.AddedLines,
					LinesRemoved: project.RemovedLines,

					Contributors: strings.Split(project.Contributors, ","),
					Pulls:        proj_pulls,
				},
			)
		}

		return OneMentorRes{Projects: mentor_stats}, 200
	} else {
		return OneMentorRes{}, 200
	}
}
