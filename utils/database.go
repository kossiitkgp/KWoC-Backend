package utils

import (
	"os"

	"kwoc20-backend/models"

	"github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/mysql" // MySQL Dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/jinzhu/gorm/dialects/sqlite" // sqlite for dev
)

var (
	DatabaseUsername = os.Getenv("DATABASE_USERNAME")
	DatabasePassword = os.Getenv("DATABASE_PASSWORD")
	DatabaseName     = os.Getenv("DATABASE_NAME")
	DatabaseHost     = os.Getenv("DATABASE_HOST")
	DatabasePort     = os.Getenv("DATABASE_PORT")
)

var newURI = "host=" + DatabaseHost + " port=" + DatabasePort + " user=" + DatabaseUsername + " dbname=" + DatabaseName + " sslmode=disable password=" + DatabasePassword

// InitialMigration Initialize migration
func InitialMigration() {
	db, err := gorm.Open("postgres", newURI)
	if err != nil {
		LOG.Println(err)
		panic(err)
	}

	db.AutoMigrate(&models.Mentor{})
	db.AutoMigrate(&models.Student{})
	db.AutoMigrate(&models.Project{})
	db.AutoMigrate(&models.PullRequest{})
	db.AutoMigrate(&models.Commits{})
}

// GetDB Get Database
func GetDB() *gorm.DB {
	db, err := gorm.Open("postgres", newURI)
	if err != nil {
		LOG.Println(err)
		panic(err)
	}

	return db
}
