package models

import (
	"github.com/jinzhu/gorm"
)

//Mentor model
type Mentor struct {
	gorm.Model

	ID 		 uint
	Name     string
	Email    string
	Username string
}

//Project Model
type Project struct {
	gorm.Model
	ID 				  uint
	Name              string
	Desc              string `gorm:"size:2550"`
	Tags              string
	RepoLink          string
	ComChannel        string
	Branch            string
	README            string `gorm:"size:255000000"`
	ProjectStatus     bool
	Mentor         	  Mentor // foreign key
	SecondaryMentor   Mentor // foreign key
}

// Commits Model
type Commits struct{
	gorm.Model
	ID uint
	URL string
	Message string
	LinesAdded uint
	LinesRemoved uint
	
	Project Project // foreign key
	Student Student // foreign key
}

// PRs Model
type PullRequest struct {
	gorm.Model
	ID	uint
	URL string
	Title string
	
	Project Project //foreign key
	Student Student // foreign key

}
