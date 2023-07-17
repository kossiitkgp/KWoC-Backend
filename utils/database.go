package utils

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/kossiitkgp/kwoc-backend/v2/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Initial database model migration
func MigrateModels(db *gorm.DB) error {
	student_mig_err := db.AutoMigrate(&models.Student{})
	if student_mig_err != nil {
		log.Err(student_mig_err).Msg("Students table automigrate error.")
		return student_mig_err
	}

	mentor_mig_err := db.AutoMigrate(&models.Mentor{})
	if mentor_mig_err != nil {
		log.Err(mentor_mig_err).Msg("Mentors table automigrate error.")
		return mentor_mig_err
	}

	project_mig_err := db.AutoMigrate(&models.Project{})
	if project_mig_err != nil {
		log.Err(project_mig_err).Msg("Projects table automigrate error.")
		return project_mig_err
	}

	stats_mig_err := db.AutoMigrate(&models.Stats{})
	if stats_mig_err != nil {
		log.Err(stats_mig_err).Msg("Stats table automigrate error.")
		return stats_mig_err
	}

	return nil
}

func GetDB() (db *gorm.DB, err error) {
	isDev := os.Getenv("DEV") == "true"

	var dialector gorm.Dialector

	if !isDev {
		DB_USERNAME := os.Getenv("DATABASE_USERNAME")
		DB_PASSWORD := os.Getenv("DATABASE_PASSWORD")
		DB_NAME := os.Getenv("DATABASE_NAME")
		DB_HOST := os.Getenv("DATABASE_HOST")
		DB_PORT := os.Getenv("DATABASE_PORT")

		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			DB_HOST,
			DB_USERNAME,
			DB_PASSWORD,
			DB_NAME,
			DB_PORT,
		)

		dialector = postgres.Open(dsn)
	} else {
		devDbPath := os.Getenv("DEV_DB_PATH")
		if devDbPath == "" {
			devDbPath = "devDB.db"
		}

		dialector = sqlite.Open(devDbPath)
	}

	db, err = gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		log.Err(err).Msg("Database open error.")

		return nil, err
	}

	return db, nil
}
