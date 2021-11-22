package controllers

import (
	"encoding/json"
	"fmt"
	"kwoc20-backend/models"
	"net/http"

	utils "kwoc20-backend/utils"
)

// ProjectReg endpoint to register project details
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
		Mentor:     mentor,
	}).Error
	if err != nil {
		utils.LOG.Println(err)
		return err.Error(), http.StatusInternalServerError
	}

	return "success", http.StatusOK
}

// ProjectGet endpoint to fetch all projects
// INCOMPLETE BECAUSE MENTOR STILL NEEDS TO BE ADDED
func AllProjects(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDB()
	defer db.Close()

	var projects []models.Project
	// Commenting Temporarily to remove Lint error as not used anywhere
	// type project_and_mentor struct {
	// 	ProjectName       string
	// 	ProjectDesc       string
	// 	ProjectTags       string
	// 	ProjectRepoLink   string
	// 	ProjectComChannel string
	// 	MentorName        []string
	// 	MentorUsername    []string
	// 	MentorEmail       []string
	// }

	err := db.Find(&projects).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// var data []project_and_mentor
	// for _, project := range projects {

	// 	mentor_names := make([]string, 1)
	// 	mentor_usernames := make([]string, 1)
	// 	mentor_emails := make([]string, 1)

	// 	var mentor models.Mentor
	// 	var secondary_mentor models.Mentor

	// 	var project_and_mentor_x project_and_mentor
	// 	err := db.First(&mentor, project.MentorID).Error
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// 	mentor_names[0] = mentor.Name
	// 	mentor_usernames[0] = mentor.Username
	// 	mentor_emails[0] = mentor.Email

	// 	if project.SecondaryMentorID != 0 {
	// 		err := db.First(&secondary_mentor, project.SecondaryMentorID).Error
	// 		if err != nil {
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			return
	// 		}
	// 		mentor_names = append(mentor_names, secondary_mentor.Name)
	// 		mentor_usernames = append(mentor_usernames, secondary_mentor.Username)
	// 		mentor_emails = append(mentor_emails, secondary_mentor.Email)
	// 	}

	// 	project_and_mentor_x.ProjectName = project.Name
	// 	project_and_mentor_x.ProjectDesc = project.Desc
	// 	project_and_mentor_x.ProjectTags = project.Tags
	// 	project_and_mentor_x.ProjectRepoLink = project.RepoLink
	// 	project_and_mentor_x.ProjectComChannel = project.ComChannel
	// 	project_and_mentor_x.MentorName = mentor_names
	// 	project_and_mentor_x.MentorUsername = mentor_usernames
	// 	project_and_mentor_x.MentorEmail = mentor_emails

	// 	data = append(data, project_and_mentor_x)
	// }
	data_json, err := json.Marshal("TO BE WORKED ON")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, error := w.Write(data_json)
	if error != nil {
		fmt.Print("ISSUE")
	}
}

// Run stats of all projects
func RunStats(req map[string]interface{}, r *http.Request) (interface{}, int) {
	test := utils.Testing()
	fmt.Println("test recieved is ", test)
	return "test", 200
}

// UpdateDetails : to Update Project Details
func UpdateDetails(req map[string]interface{}, r *http.Request) (interface{}, int) {
	db := utils.GetDB()
	defer db.Close()

	project := &models.Project{
		Desc:   req["desc"].(string),
		Tags:   req["tags"].(string),
		Branch: req["branch"].(string),
	}
	fmt.Print(project)
	err := db.Table("projects").Where("Name= ?", req["name"].(string)).Select("Desc", "Tags", "Branch").Updates(project).Error
	if err != nil {
		fmt.Print(err)
		return "fail", http.StatusBadRequest
	}

	return "Success", http.StatusOK
}

// ProjectDetails fetch endpoint
func ProjectDetails(req map[string]interface{}, r *http.Request) (interface{}, int) {
	db := utils.GetDB()
	defer db.Close()

	name := req["name"].(string)
	fmt.Print(name)
	projects := models.Project{}
	err := db.Where(&models.Project{Name: name}).First(&projects).Error
	if err != nil {
		return err, http.StatusBadRequest
	}

	type Response map[string]interface{}
	res := Response{
		"name":   projects.Name,
		"desc":   projects.Desc,
		"tags":   projects.Tags,
		"branch": projects.Branch,
	}
	fmt.Print(projects.Tags)
	return res, http.StatusOK
}
