package models

import "gorm.io/gorm"

// `projects` table models

// Projects table row
type Project struct {
	gorm.Model

	Name          string
	Desc          string `gorm:"size:2550"`
	Tags          string
	RepoLink      string
	ComChannel    string
	README        string `gorm:"size:2550000000"`
	ProjectStatus bool   `gorm:":default:false"`

	// for stats
	LastPullTime int64

	// stats table
	CommitCount  uint
	PRCount      uint
	AddedLines   uint
	RemovedLines uint

	// list of students who contributed to the project (a string of usernames separated by comma(,))
	Contributors string

	// list of URLs to PRs contributed to the project (a string of links separated by comma(,))
	Pulls string

	// foreign keys
	Mentor_id          int32
	Mentor             Mentor `gorm:"ForeignKey:Mentor_id"`
	SecondaryMentor_id int32
	SecondaryMentor    Mentor `gorm:"ForeignKey:SecondaryMentor_id"`
}
