package controllers

import (
	"fmt"
	"log"
	"net/http"

	"kwoc20-backend/models"

	"kwoc20-backend/utils"

	"github.com/gorilla/mux"
)

type Student struct {
	ID           uint
	Username     string
	CommitCount  uint
	PRCount      uint
	AddedLines   uint
	RemovedLines uint
	TechWorked   string
}

func CheckStudent(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	params := mux.Vars(r)
	var student Student
	db.
	Table("students").
	Select(
		"id", "username", "commit_count",
		"pr_count", "added_lines", "removed_lines",
		"tech_worked",
	).
	Find(&student, params["username"])

	if student.Username == " " {
		w.WriteHeader(400)
		_, err := w.Write([]byte("false"))
		if err != nil {
			log.Printf("Write failed: %v", err)
		}
		return
	} else {
		w.WriteHeader(200)
		_, err := w.Write([]byte("true"))
		if err != nil {
			log.Printf("Write failed: %v", err)
		}
		return
	}
}

func AllStudents(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	w.WriteHeader(200)
	var students []Student
	db.
	Table("students").
	Select(
		"id", "username", "commit_count",
		"pr_count", "added_lines", "removed_lines",
		"tech_worked",
	).
	Find(&students)
	str := fmt.Sprintf("%+v", students)
	_, err := w.Write([]byte(str))
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
	return
}

func OneStudent(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	w.WriteHeader(200)
	params := mux.Vars(r)
	var student Student
	db.
	Table("students").
	Select(
		"id", "username", "commit_count",
		"pr_count", "added_lines", "removed_lines",
		"tech_worked",
	).
	Find(&student, params["username"])
	str := fmt.Sprintf("%+v", student)
	_, err := w.Write([]byte(str))
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
	return
}

type Project struct {
	ID uint
	RepoLink      string
	Branch        string

	LastPullDate  string

	CommitCount   uint
	PRCount       uint
	AddedLines    uint
	RemovedLines  uint
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
		log.Printf("Write failed: %v", err)
	}
	return
}

func OneMentor(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	w.WriteHeader(200)
	params := mux.Vars(r)
	var mentor []models.Project
	db.
		Table("projects").
		Where("project_status = ?", "true").
		Select(
			"id", "repo_link", "branch",
			"last_pull_date", "commit_count",
			"pr_count", "added_lines", "removed_lines",
		).
		Find(&mentor, params["Mentor.username"])
	str := fmt.Sprintf("%+v", mentor)
	_, err := w.Write([]byte(str))
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
	return
}
