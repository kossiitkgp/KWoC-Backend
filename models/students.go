package models

// `students` table model

// Students table row
type Student struct {
	Model

	Name           string `gorm:"column:name"`
	Email          string `gorm:"column:email"`
	College        string `gorm:"column:college"`
	Username       string `gorm:"column:username"`
	PassedMidEvals bool   `gorm:"column:passed_mid_evals"`
	PassedEndEvals bool   `gorm:"column:passed_end_evals"`
	BlogLink       string `gorm:"size:2550;column:blog_link"`

	// stats table
	CommitCount  uint `gorm:"column:commit_count"`
	PullCount    uint `gorm:"column:pull_count"`
	LinesAdded   uint `gorm:"column:lines_added"`
	LinesRemoved uint `gorm:"column:lines_removed"`

	// LanguagesUsed is a string of languages separated by comma(,)
	LanguagesUsed string `gorm:"column:languages_used"`

	// ProjectsWorked is a string of project IDs separated by comma(,)
	ProjectsWorked string `gorm:"column:projects_worked"`

	// list of URLs to PRs contributed by the student (a string of links separated by comma(,))
	Pulls string `gorm:"column:pulls"`
}
