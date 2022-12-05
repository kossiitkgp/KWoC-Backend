package controllers

import (
	"net/http"

	"kwoc20-backend/models"

	"kwoc20-backend/utils"

	"github.com/gorilla/mux"
)

func CheckStudent(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()

	params := mux.Vars(r)
	var student models.Student
	db.Find(&student, params["Username"])
	if student.Username == " " {
		w.WriteHeader(400)
		w.Write([]byte("false"))
	} else {
		w.WriteHeader(200)
		w.Write([]byte("true"))
	}
}

func AllStudents(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	w.WriteHeader(200)
	var students []models.Student
	db.Find(&students)
	w.Write(students)
}

func OneStudent(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	w.WriteHeader(200)
	params := mux.Vars(r)
	var student models.Student
	db.Find(&student, params["Username"])
	w.Write(student)
}

func GetAllProjects(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	w.WriteHeader(200)
	var projects []models.Project
	db.Find(&projects)
	w.Write(projects)
}

func OneMentor(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	w.WriteHeader(200)
	params := mux.Vars(r)
	var mentor []models.Project
	db.Find(&mentor, params["Mentor.Username"])
	w.Write(mentor)
}
