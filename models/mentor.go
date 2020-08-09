package models

import (
	"github.com/jinzhu/gorm"
)

//Mentor model
type Mentor struct {
	gorm.Model
	
	Name string
	Email string
	Username string
}
