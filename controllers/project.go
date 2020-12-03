package controllers

import (
	"encoding/json"
	"net/http"

	"kwoc20-backend/models"
	utils "kwoc20-backend/utils"
)

//ProjectReg endpoint to register project details
func ProjectReg(req map[string]interface{}, r *http.Request) (interface{}, int) {

	db := utils.GetDB()
	defer db.Close()

	gh_username := req["username"].(string)

	ctx_user := r.Context().Value(utils.CtxUserString("user")).(string)

	if ctx_user != gh_username {
		utils.LOG.Printf("%v != %v Detected Session Hijacking\n", gh_username, ctx_user)
		return "Corrupt JWT", http.StatusForbidden
	}

	mentor := models.Mentor{}
	db.Where(&models.Mentor{Username: gh_username}).First(&mentor)

	err := db.Create(&models.Project{
		Name:       req["name"].(string),
		Desc:       req["desc"].(string),
		Tags:       req["tags"].(string),
		RepoLink:   req["repoLink"].(string),
		ComChannel: req["comChannel"].(string),
		MentorID:   mentor.ID,
	}).Error

	if err != nil {
		utils.LOG.Println(err)
		return err.Error(), http.StatusInternalServerError
	}

	return "success", http.StatusOK

}

//ProjectGet endpoint to fetch all projects
// INCOMPLETE BECAUSE MENTOR STILL NEEDS TO BE ADDED
func ProjectGet(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	defer db.Close()

	var projects []models.Project
	err := db.Find(&projects).Error
	if err != nil {
		http.Error(w, err.Error(), 400)
		utils.LOG.Println(err)
		return
	}

	err = json.NewEncoder(w).Encode(projects)
	if err != nil {
		http.Error(w, err.Error(), 500)
		utils.LOG.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`success`))

}
