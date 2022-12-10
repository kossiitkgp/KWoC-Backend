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

type AllProjectsProject struct {
	Title   string
	Link    string
	Contri  uint // Number of students who contributed
	Commits uint
	Lines   string
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
				Title:   project.Name,
				Link:    project.RepoLink,
				Contri:  0,
				Commits: project.CommitCount,
				Lines:   fmt.Sprintf("+%d/-%d", project.AddedLines, project.RemovedLines),
			},
		)
	}

	return AllProjectsRes{Stats: response}, 200
}

type OneMentorProj struct {
	RepoLink     string
	Project_name string
	Contributors []string // Array of usernames of students who contributed
	Commits      uint
	LinesAdded   uint
	LinesRemoved uint
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
			mentor_stats = append(
				mentor_stats,
				OneMentorProj{
					RepoLink:     project.RepoLink,
					Project_name: project.Name,
					Commits:      project.CommitCount,
					LinesAdded:   project.AddedLines,
					LinesRemoved: project.RemovedLines,
				},
			)
		}

		return OneMentorRes{Projects: mentor_stats}, 200
	} else {
		return OneMentorRes{}, 200
	}
}
