package models

import "gorm.io/gorm"

// `stats` table model

// Stats table row
type Stats struct {
	gorm.Model

	TotalCommits      uint
	TotalPRCount      uint
	TotalLinesAdded   uint
	TotalLinesRemoved uint

	// Time at which the stats in this entry were generated
	GenTime int64
}
