package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"fmt"

	"kwoc20-backend/models"
	logs "kwoc20-backend/utils/logs/pkg"

	"github.com/go-kit/kit/log/level"
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
		level.Error(logs.Logger).Log("error", fmt.Sprintf("%v", err))
		return
	}

	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		level.Error(logs.Logger).Log("error", fmt.Sprintf("%v", err))
		return
	}
	defer db.Close()

	err = db.Create(&models.Mentor{
		Name:         mentor.Name,
		Email:        mentor.Email,
		GithubHandle: mentor.GithubHandle,
		AccessToken:  mentor.AccessToken,
	}).Error

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		level.Error(logs.Logger).Log("error", fmt.Sprintf("%v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message" : "success"}`))

}
