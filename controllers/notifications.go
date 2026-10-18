package controllers

import (
	"fmt"
	"net/http"

	"github.com/kossiitkgp/kwoc-backend/v2/models"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"

	"gorm.io/gorm"
)

// projectMentors returns the primary and secondary mentors of a project.
func projectMentors(db *gorm.DB, project models.Project) []models.Mentor {
	mentors := make([]models.Mentor, 0, 2)
	seen := make(map[uint]bool)

	addMentor := func(id uint) {
		if id == 0 || seen[id] {
			return
		}

		mentor := models.Mentor{}
		if err := db.Where("id = ?", id).First(&mentor).Error; err != nil {
			return
		}

		seen[id] = true
		mentors = append(mentors, mentor)
	}

	addMentor(uint(project.MentorId))
	if project.SecondaryMentorId != nil {
		addMentor(uint(*project.SecondaryMentorId))
	}

	return mentors
}

// notifyProjectSubmitted tells the mentors that their project is awaiting
// admin approval.
func notifyProjectSubmitted(r *http.Request, db *gorm.DB, project models.Project) {
	if !utils.IsEmailConfigured() {
		return
	}

	actionURL := utils.FrontendURL()
	if actionURL != "" {
		actionURL += "/mentor/dashboard"
	}

	heading := "Project submitted for approval"
	message := fmt.Sprintf(
		"The project \"%s\" was submitted successfully and is waiting for admin approval. You will receive another email once it has been reviewed.",
		project.Name,
	)

	for _, mentor := range projectMentors(db, project) {
		if mentor.Email == "" {
			continue
		}

		if err := utils.SendActionRequiredEmail(
			mentor.Email,
			mentor.Name,
			mentor.Username,
			heading,
			message,
			actionURL,
			"View project",
		); err != nil {
			utils.LogErr(r, err, fmt.Sprintf("Error sending project submission email to `%s`.", mentor.Username))
		}
	}
}

// notifyProjectApproved tells the mentors that their project is live.
func notifyProjectApproved(r *http.Request, db *gorm.DB, project models.Project) {
	if !utils.IsEmailConfigured() {
		return
	}

	for _, mentor := range projectMentors(db, project) {
		if mentor.Email == "" {
			continue
		}

		if err := utils.SendProjectApprovedEmail(
			mentor.Email,
			mentor.Name,
			mentor.Username,
			project.Name,
		); err != nil {
			utils.LogErr(r, err, fmt.Sprintf("Error sending project approval email to `%s`.", mentor.Username))
		}
	}
}
