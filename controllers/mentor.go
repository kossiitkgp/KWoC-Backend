package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kwoc20-backend/models"
	"net/http"

	"github.com/jinzhu/gorm"
)

// MentorOauth Handler for Github OAuth of Mentor
func MentorOAuth(w http.ResponseWriter, r *http.Request) {
	// get the code from frontend
	type MentorOAuthCode struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}
	var mentorOAuthCode MentorOAuthCode
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &mentorOAuthCode)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	// using the code obtained from above to get AccessToken from Github
	req, _ := json.Marshal(map[string]string{
		"client_id":     "74557dcb91016b10b54b",
		"client_secret": "594d9e729a47a5d8e944edc792530342841caaf0",
		"code":          mentorOAuthCode.Code,
		"state":         "PAKKA RANDOM",
	})
	res, err := http.Post("https://github.com/login/oauth/access_token", "application/json", bytes.NewBuffer(req))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	resBody, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(resBody))

	w.Write([]byte("helo"))
}

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
