package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"kwoc20-backend/models"

	"github.com/jinzhu/gorm"
)

//MentorReg Handler for Registering Mentors
func MentorReg(w http.ResponseWriter, r *http.Request) {

	var mentor models.Mentor
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &mentor)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	defer db.Close()

	db.Create(&models.Mentor{
		Name:         mentor.Name,
		Email:        mentor.Email,
		GithubHandle: mentor.GithubHandle,
		AccessToken:  mentor.AccessToken,
	})

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"message": "` + err.Error() + `"}`))
	return

}
