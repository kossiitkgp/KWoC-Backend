package controllers

import (
	"fmt"
	"kwoc20-backend/models"
	"net/http"

	utils "kwoc20-backend/utils"
)

// ProjectReg endpoint to register project details
func ProjectReg(req map[string]interface{}, r *http.Request) (interface{}, int) {
	/*
		 BODY PARAMS
			{
				"id" : Project Id,
				"name" : New Name of Project,
				"desc" : New DEsciption of Project,
				"tags" : Updated tags of project,
				"branch" : updated branch,
				"username" : Mentor Username,
				"secondaryMentor": Secodnary Mentor useraName
				"repoLink": RepoLink of Project,
				"comChannel" : Link of communication channel of mentor and mentee,
				"readme" : Project Readme
			}
	*/
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

	secondaryMentor := models.Mentor{}
	db.Where(&models.Mentor{Username: req["secondaryMentor"].(string)}).First(&secondaryMentor)

	err := db.Create(&models.Project{
		Name:            req["name"].(string),
		Desc:            req["desc"].(string),
		Tags:            req["tags"].(string),
		RepoLink:        req["repoLink"].(string),
		ComChannel:      req["comChannel"].(string),
		README:          req["readme"].(string),
		Branch:          req["branch"].(string),
		Mentor:          mentor,
		SecondaryMentor: secondaryMentor,
	}).Error
	if err != nil {
		utils.LOG.Println(err)
		return err.Error(), http.StatusInternalServerError
	}

	return "success", http.StatusOK
}

// ProjectGet endpoint to fetch all projects
// INCOMPLETE BECAUSE MENTOR STILL NEEDS TO BE ADDED
func AllProjects(req map[string]interface{}, r *http.Request) (interface{}, int) {
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

	err := db.Not("project_status", false).Find(&projects).Error
	if err != nil {
		fmt.Print(err)
		return "fail", http.StatusInternalServerError
	}

	return projects, 200

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

}

// Run stats of all projects
func RunStats(req map[string]interface{}, r *http.Request) (interface{}, int) {
	test := utils.Testing()
	fmt.Println("test recieved is ", test)
	return "test", 200
}

// UpdateDetails : to Update Project Details
func UpdateDetails(req map[string]interface{}, r *http.Request) (interface{}, int) {
	/*
		 BODY PARAMS
			{
				"id" : Project Id,
				"name" : New Name of Project,
				"desc" : New DEsciption of Project,
				"tags" : Updated tags of project,
				"branch" : updated branch,
				"readme" :  Project Readme
			}
	*/
	db := utils.GetDB()
	defer db.Close()

	ctx_user := r.Context().Value(utils.CtxUserString("user")).(string)

	secondaryMentor := models.Mentor{}
	db.Where(&models.Mentor{Username: req["secondaryMentor"].(string)}).First(&secondaryMentor)

	id := (uint)(req["id"].(float64))
	project := &models.Project{
		Name:            req["name"].(string),
		Desc:            req["desc"].(string),
		Tags:            req["tags"].(string),
		Branch:          req["branch"].(string),
		README:          req["readme"].(string),
		SecondaryMentor: secondaryMentor,
	}
	fmt.Print(project)
	projects := models.Project{}
	err := db.First(&projects, id).Select("Name", "Desc", "Tags", "Branch", "README", "SecondaryMentor_id").Updates(project).Error
	if err != nil {
		fmt.Print(err)
		return "fail", http.StatusBadRequest
	}

	if projects.SecondaryMentor.Username != ctx_user {
		fmt.Println(projects.Mentor.Username, "+", "ctx_user")
		return "Session Hijacking", 403
	}

	return "Success", http.StatusOK
}

// ProjectDetails fetch endpoint
func ProjectDetails(req map[string]interface{}, r *http.Request) (interface{}, int) {
	/*
		 BODY PARAMS
			{
				"id" : Project Id,
			}
	*/
	db := utils.GetDB()
	defer db.Close()

	ctx_user := r.Context().Value(utils.CtxUserString("user")).(string)

	id := (uint)(req["id"].(float64))

	projects := models.Project{}

	err := db.First(&projects, id).Error
	if err != nil {
		return err, http.StatusBadRequest
	}

	if projects.Mentor.Username != ctx_user {
		fmt.Println(projects.Mentor.Username, "+", "ctx_user")
		return "Session Hijacking", 403
	}

	type Response map[string]interface{}
	res := Response{
		"name":            projects.Name,
		"desc":            projects.Desc,
		"tags":            projects.Tags,
		"branch":          projects.Branch,
		"repo_link":       projects.RepoLink,
		"secondaryMentor": projects.SecondaryMentor.Username,
	}
	return res, http.StatusOK
}
