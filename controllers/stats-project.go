package controllers

import (
	"kwoc20-backend/models"
	"kwoc20-backend/utils"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type AllProjectsProject struct {
	Name string `json:"name"`
	Link string `json:"link"`

	CommitCount  uint `json:"commit_count"`
	PrCount      uint `json:"pr_count"`
	LinesAdded   uint `json:"lines_added"`
	LinesRemoved uint `json:"lines_removed"`
}
type AllProjectsRes struct {
	Stats []AllProjectsProject
}

func GetAllProjects(req map[string]interface{}, r *http.Request) (interface{}, int) {
	db := utils.GetDB()
	defer db.Close()

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
				PrCount:      project.PRCount,
				LinesAdded:   project.AddedLines,
				LinesRemoved: project.RemovedLines,
			},
		)
	}

	return AllProjectsRes{Stats: response}, 200
}

type OneMentorProjPull struct {
	Url string `json:"url"`
}
type OneMentorProj struct {
	Name     string `json:"name"`
	RepoLink string `json:"repo_link"`

	CommitCount  uint `json:"commit_count"`
	LinesAdded   uint `json:"lines_added"`
	LinesRemoved uint `json:"lines_removed"`

	Contributors []string            `json:"contributors"` // Array of usernames of students who contributed
	Pulls        []OneMentorProjPull `json:"pulls"`
}
type OneMentorRes struct {
	Projects []OneMentorProj
}

func OneMentor(req map[string]interface{}, r *http.Request) (interface{}, int) {
	db := utils.GetDB()
	defer db.Close()

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
		return "Error: Mentor Does Not Exist", 410
	}
}
