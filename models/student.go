package models

import (
	"github.com/jinzhu/gorm"
)

// Mentor model
type Student struct {
	gorm.Model

	Name     string
	Email    string
	College  string
	Username string
	Evals    int    `gorm:"default:0"`
	BlogLink string `gorm:"size:2550"`

	// stats table
	CommitCount  uint
	PRCount      uint
	AddedLines   uint
	RemovedLines uint

	// TechWorked is a string of languages separated by comma(,)
	TechWorked string

	// ProjectsWorked is a string of project IDs separated by comma(,)
	ProjectsWorked string
}
