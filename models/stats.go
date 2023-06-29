package models

import "gorm.io/gorm"

// `stats` table model

// Stats table row
type Stats struct {
	gorm.Model

	TotalCommitCount  uint `gorm:"column:total_commit_count"`
	TotalPullCount    uint `gorm:"column:total_pull_count"`
	TotalLinesAdded   uint `gorm:"column:total_lines_added"`
	TotalLinesRemoved uint `gorm:"column:total_lines_removed"`

	// Time at which the stats in this entry were generated
	GenTime int64
}
