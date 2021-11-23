package models

import (
	"github.com/jinzhu/gorm"
)

// Mentor model
type Mentor struct {
	gorm.Model

	Name     string
	Email    string
	Username string
}

// Project Model
type Project struct {
	gorm.Model

	Name          string
	Desc          string `gorm:"size:2550"`
	Tags          string
	RepoLink      string
	ComChannel    string
	Branch        string
	README        string `gorm:"size:2550000000"`
	ProjectStatus bool

	// for stats
	LastCommitSHA string
	LastPullDate  string

	// stats table
	CommitCount  uint
	PRCount      uint
	AddedLines   uint
	RemovedLines uint

	// foreign keys
	MentorUsername          string
	Mentor                  Mentor `gorm:"foreignKey:MentorUsername"`
	SecondaryMentorUsername string
	SecondaryMentor         Mentor `gorm:"foreignKey:SecondaryMentorUsername"`
}

// Commits Model
type Commits struct {
	gorm.Model

	URL          string
	Message      string
	LinesAdded   uint
	LinesRemoved uint
	SHA          string

	Project Project // foreign key
	Student Student // foreign key
}

// PRs Model
type PullRequest struct {
	gorm.Model

	URL       string
	Title     string
	CreatedAt string

	Project Project // foreign key
	Student Student // foreign key

}
