package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func AllStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var students []Student
	DB.Find(&students)
	json.NewEncoder(w).Encode(students)
}

func OneStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var student Student
	DB.Find(&student, params["username"])
	json.NewEncoder(w).Encode(student)
}

func AllProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var projects []Project
	DB.Find(&projects)
	json.NewEncoder(w).Encode(projects)
}
