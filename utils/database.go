package utils

import (
	"kwoc20-backend/models"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"

	// _ "github.com/jinzhu/gorm/dialects/mysql" // MySQL Dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // sqlite for dev
)

// InitialMigration Initialize migration
func InitialMigration() {
	db := GetDB()
	defer db.Close()

	db.AutoMigrate(&models.Mentor{})
	db.AutoMigrate(&models.Student{})
	db.AutoMigrate(&models.Project{})
	db.AutoMigrate(&models.PullRequest{})
	db.AutoMigrate(&models.Commits{})
}

// GetDB Get Database
func GetDB() *gorm.DB {
	isDev := os.Getenv("DEV") == "true"

	if !isDev {
		DatabaseUsername := os.Getenv("DATABASE_USERNAME")
		DatabasePassword := os.Getenv("DATABASE_PASSWORD")
		DatabaseName := os.Getenv("DATABASE_NAME")
		DatabaseHost := os.Getenv("DATABASE_HOST")
		DatabasePort := os.Getenv("DATABASE_PORT")

		newURI := "host=" + DatabaseHost + " port=" + DatabasePort + " user=" + DatabaseUsername + " dbname=" + DatabaseName + " sslmode=disable password=" + DatabasePassword
		db, err := gorm.Open("postgres", newURI)

		if err != nil {
			log.Err(err).Msg("Database Error")
			panic(err)
		}

		return db
	} else {
		// SQLite database for local development

		db, err := gorm.Open("sqlite3", "devDB.db")
		if err != nil {
			log.Err(err).Msg("Database Error")
			panic(err)
		}

		return db
	}
}
