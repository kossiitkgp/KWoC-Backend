package models

import "gorm.io/gorm"

// `mentors` table models

// Mentors table row
type Mentor struct {
	gorm.Model

	Name     string
	Email    string
	Username string
}
