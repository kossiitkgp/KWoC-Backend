package controllers

import (
	"fmt"
	"net/http"

	"kwoc20-backend/models"

	"kwoc20-backend/utils"

	"github.com/gorilla/mux"
)

func CheckStudent(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()

	params := mux.Vars(r)
	var student models.Student
	db.Find(&student, params["username"])
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
	str := fmt.Sprintf("%+v", students)
	w.Write([]byte(str))
}

func OneStudent(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	w.WriteHeader(200)
	params := mux.Vars(r)
	var student models.Student
	db.Find(&student, params["username"])
	str := fmt.Sprintf("%+v", student)
	w.Write([]byte(str))
}

func GetAllProjects(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	w.WriteHeader(200)
	var projects []models.Project
	db.Find(&projects)
	str := fmt.Sprintf("%+v", projects)
	w.Write([]byte(str))
}

func OneMentor(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	w.WriteHeader(200)
	params := mux.Vars(r)
	var mentor []models.Project
	db.Find(&mentor, params["Mentor.username"])
	str := fmt.Sprintf("%+v", mentor)
	w.Write([]byte(str))
}
