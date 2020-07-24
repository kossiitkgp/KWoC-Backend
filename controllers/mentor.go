package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	
	"kwoc20-backend/models"
	utils "kwoc20-backend/utils"
	
	"github.com/jinzhu/gorm"
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

	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		http.Error(w, err.Error(), 500)
		utils.LOG.Println(err)
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
		http.Error(w, err.Error(), 500)
		utils.LOG.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`success`))

}
