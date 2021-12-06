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
	ProjectStatus bool   `gorm:":default:false"`

	// for stats
	LastCommitSHA string
	LastPullDate  string

	// stats table
	CommitCount  uint
	PRCount      uint `gorm:":default:0"`
	AddedLines   uint
	RemovedLines uint

	// foreign keys
	Mentor_id          int32
	Mentor             Mentor `gorm:"ForeignKey:Mentor_id"`
	SecondaryMentor_id int32
	SecondaryMentor    Mentor `gorm:"ForeignKey:SecondaryMentor_id"`
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

	Project Project // foreign key
	Student Student // foreign key

}
