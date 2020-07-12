package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"fmt"

	"kwoc20-backend/models"
	logs "kwoc20-backend/utils/logs/pkg"

	"github.com/jinzhu/gorm"
	"github.com/go-kit/kit/log/level"
)

//MentorReg Handler for Registering Mentors
func MentorReg(w http.ResponseWriter, r *http.Request) {

	var mentor models.Mentor
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, writeErr := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		if writeErr !=nil {
			_ = level.Warn(logs.Logger).Log("error",fmt.Sprintf("%v",writeErr))
		}
		_ = level.Error(logs.Logger).Log("error", fmt.Sprintf("%v",err))
		return
	}
	err = json.Unmarshal(body, &mentor)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, writeErr := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		if writeErr !=nil {
			_ = level.Warn(logs.Logger).Log("error",fmt.Sprintf("%v",writeErr))
		}
		_ = level.Error(logs.Logger).Log("error", fmt.Sprintf("%v",err))
		return
	}

	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, writeErr := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		if writeErr !=nil {
			_ = level.Warn(logs.Logger).Log("error",fmt.Sprintf("%v",writeErr))
		}
		_ = level.Error(logs.Logger).Log("error", fmt.Sprintf("%v",err))
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
	_, writeErr := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
	if writeErr !=nil {
		_ = level.Warn(logs.Logger).Log("error",fmt.Sprintf("%v",writeErr))
	}
	_ = level.Error(logs.Logger).Log("error", fmt.Sprintf("%v",err))
}
