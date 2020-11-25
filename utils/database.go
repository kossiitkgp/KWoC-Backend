package utils

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // MySQL Dialect

	"kwoc20-backend/models"
)

// InitialMigration Initialize migration
func InitialMigration() {
	DatabaseUsername := os.Getenv("DATABASE_USERNAME")
	DatabasePassword := os.Getenv("DATABASE_PASSWORD")
	DatabaseName := os.Getenv("DATABASE_NAME")
	DatabaseHost := os.Getenv("DATABASE_HOST")
	DatabasePort := os.Getenv("DATABASE_PORT")

	DatabaseURI := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		DatabaseUsername,
		DatabasePassword,
		DatabaseHost,
		DatabasePort,
		DatabaseName,
	)

	db, err := gorm.Open("mysql", DatabaseURI)
	if err != nil {
		LOG.Println(err)
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&models.Mentor{})
	db.AutoMigrate(&models.Project{})
	db.AutoMigrate(&models.Student{})
}

// GetDB Get Database
func GetDB() *gorm.DB {
	DatabaseUsername := os.Getenv("DATABASE_USERNAME")
	DatabasePassword := os.Getenv("DATABASE_PASSWORD")
	DatabaseName := os.Getenv("DATABASE_NAME")
	DatabaseHost := os.Getenv("DATABASE_HOST")
	DatabasePort := os.Getenv("DATABASE_PORT")

	DatabaseURI := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		DatabaseUsername,
		DatabasePassword,
		DatabaseHost,
		DatabasePort,
		DatabaseName,
	)

	db, err := gorm.Open("mysql", DatabaseURI)
	if err != nil {
		LOG.Println(err)
		panic(err)
	}

	return db
}
