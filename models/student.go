package models

import (
	"github.com/jinzhu/gorm"
)

//Mentor model
type Student struct {
	gorm.Model

	Name     string
	Email    string
	College  string
	Username string
	Evals    int 	`gorm:"default:0"`
}
