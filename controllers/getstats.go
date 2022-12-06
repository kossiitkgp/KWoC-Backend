package controllers

import (
	"fmt"
	"log"
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
	var students []models.Student
	db.Find(&students)
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
	var student models.Student
	db.Find(&student, params["username"])
	str := fmt.Sprintf("%+v", student)
	_, err := w.Write([]byte(str))
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
	return
}

func GetAllProjects(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	w.WriteHeader(200)
	var projects []models.Project
	db.Find(&projects)
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
	db.Find(&mentor, params["Mentor.username"])
	str := fmt.Sprintf("%+v", mentor)
	_, err := w.Write([]byte(str))
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
	return
}
