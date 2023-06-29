package models

import "gorm.io/gorm"

// `mentors` table models

// Mentors table row
type Mentor struct {
	gorm.Model

	Name     string `gorm:"column:name"`
	Email    string `gorm:"column:email"`
	Username string `gorm:"column:username"`
}
