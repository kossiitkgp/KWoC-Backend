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
		return err.Error(), 500
	}

	return "success", 200

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

func AllProjects(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	defer db.Close()

	var projects []models.Project
	type project_and_mentor struct {
		ProjectName       		string
		ProjectDesc       		string
		ProjectTags       		string
		ProjectRepoLink   		string
		ProjectComChannel 		string
		MentorName 				string	
		MentorUsername 			string
		MentorEmail				string
	}

	err := db.Find(&projects).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var data []project_and_mentor
	for _, project := range projects {
		var mentor models.Mentor
		var project_and_mentor_x project_and_mentor
		err := db.First(&mentor, project.MentorID).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		project_and_mentor_x.ProjectName = 			project.Name
		project_and_mentor_x.ProjectDesc = 			project.Desc
		project_and_mentor_x.ProjectTags = 			project.Tags
		project_and_mentor_x.ProjectRepoLink = 		project.RepoLink
		project_and_mentor_x.ProjectComChannel = 	project.ComChannel
		project_and_mentor_x.MentorName = 			mentor.Name
		project_and_mentor_x.MentorUsername = 		mentor.Username
		project_and_mentor_x.MentorEmail = 			mentor.Email

		data = append(data, project_and_mentor_x)	
	}
	data_json, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data_json)
}
