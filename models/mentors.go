package models

// `mentors` table models

// Mentors table row
type Mentor struct {
	Model

	Name     string `gorm:"column:name"`
	Email    string `gorm:"column:email"`
	Username string `gorm:"column:username"`
}
