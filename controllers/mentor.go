package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"kwoc20-backend/models"
	"kwoc20-backend/utils"
)

//MentorReg Handler for Registering Mentors
func MentorReg(w http.ResponseWriter, r *http.Request) {
	var mentor models.Mentor
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &mentor)
	if err != nil {
		http.Error(w, err.Error(), 400)
		utils.LOG.Println(err)
		return
	}

	db := utils.GetDB()
	defer db.Close()

	err = db.Create(&models.Mentor{
		Name:         mentor.Name,
		Email:        mentor.Email,
		GithubHandle: mentor.GithubHandle,
		AccessToken:  mentor.AccessToken,
	}).Error

	if err != nil {
		http.Error(w, err.Error(), 500)
		utils.LOG.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`success`))

}
