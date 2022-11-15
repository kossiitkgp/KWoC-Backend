package controllers

import (
	"fmt"
	"kwoc20-backend/models"
	"net/http"

	utils "kwoc20-backend/utils"

	"github.com/rs/zerolog/log"
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
		log.Warn().Msgf(
			"%s %s: %v != %v Detected Session Hijacking",
			r.Method,
			r.RequestURI,
			gh_username,
			ctx_user,
		)
		return "Corrupt JWT", http.StatusForbidden
	}

	mentor := models.Mentor{}
	db.Where(&models.Mentor{Username: gh_username}).First(&mentor)

	secondaryMentor := models.Mentor{}
	if len(req["secondaryMentor"].(string)) > 0 {
		db.Where(&models.Mentor{Username: req["secondaryMentor"].(string)}).First(&secondaryMentor)
	}

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
		utils.LogErr(
			r,
			err,
			"Database error.",
		)
		return err.Error(), http.StatusInternalServerError
	}

	return "success", http.StatusOK
}

// ProjectGet endpoint to fetch all projects
func AllProjects(req map[string]interface{}, r *http.Request) (interface{}, int) {
	db := utils.GetDB()
	defer db.Close()

	var projects []models.Project

	err := db.Preload("Mentor").Preload("SecondaryMentor").Not("project_status", "false").Find(&projects).Error
	if err != nil {
		utils.LogErr(
			r,
			err,
			"Database Error",
		)
		return "fail", http.StatusInternalServerError
	}

	return projects, 200
}

// Run stats of all projects
func RunStats(req map[string]interface{}, r *http.Request) (interface{}, int) {
	test := utils.Testing()
	utils.LogInfo(
		r,
		fmt.Sprintf("test recieved is ", test),
	)
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
				"readme" :  Project Readme,
				"secondaryMentor":Secondary Mentor Username,
				"comChannel": Communication Channel
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
		ComChannel:      req["comChannel"].(string),
	}
	projects := models.Project{}
	err := db.Preload("Mentor").First(&projects, id).Select("Name", "Desc", "Tags", "Branch", "README", "SecondaryMentor", "ComChannel").Updates(project).Error
	if err != nil {
		utils.LogWarn(
			r,
			fmt.Sprintf("Bad Request %v", err),
		)
		return "fail", http.StatusBadRequest
	}

	if projects.Mentor.Username != ctx_user {
		utils.LogWarn(
			r,
			fmt.Sprintf(
				"%v != %v Detected Session Hijacking",
				projects.Mentor.Username,
				ctx_user,
			),
		)
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

	err := db.Preload("Mentor").Preload("SecondaryMentor").First(&projects, id).Error
	if err != nil {
		return err, http.StatusBadRequest
	}
	if projects.Mentor.Username != ctx_user {
		utils.LogWarn(
			r,
			fmt.Sprintf(
				"%v != %v Detected Session Hijacking",
				projects.Mentor.Username,
				ctx_user,
			),
		)
		return "Session Hijacking", 403
	}

	type Response map[string]interface{}
	res := Response{
		"name":            projects.Name,
		"desc":            projects.Desc,
		"tags":            projects.Tags,
		"branch":          projects.Branch,
		"repo_link":       projects.RepoLink,
		"comChannel":      projects.ComChannel,
		"secondaryMentor": projects.SecondaryMentor.Username,
	}
	return res, http.StatusOK
}
