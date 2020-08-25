package models

import (
	"github.com/jinzhu/gorm"
)

//Mentor model
type Student struct {
	gorm.Model
	
	Name string
	Email string
	College string	
	Username string

	MidsCleared bool
	EndsCleared bool
}
