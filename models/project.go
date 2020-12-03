package models

import (
	"github.com/jinzhu/gorm"
)

//Project Model
type Project struct {
	gorm.Model
	Name       string
	Desc       string `gorm:"size:2550"`
	Tags       string
	RepoLink   string
	ComChannel string
	MentorID   uint
	SecondaryMentorID uint
}
