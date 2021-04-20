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
	Name              string
	Desc              string `gorm:"size:2550"`
	Tags              string
	RepoLink          string
	ComChannel        string
	Branch            string
	README            string `gorm:"size:255000000"`
	ProjectStatus     bool
	Mentor         	  Mentor
	SecondaryMentor   Mentor
}
