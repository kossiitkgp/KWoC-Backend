package utils

import (
	"fmt"
	// "os"

	"github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/mysql" // MySQL Dialect
	_ "github.com/jinzhu/gorm/dialects/sqlite" //sqlite for dev

	"kwoc20-backend/models"
)

// InitialMigration Initialize migration
func InitialMigration() {
	// DatabaseUsername := os.Getenv("DATABASE_USERNAME")
	// DatabasePassword := os.Getenv("DATABASE_PASSWORD")
	// DatabaseName := os.Getenv("DATABASE_NAME")
	// DatabaseHost := os.Getenv("DATABASE_HOST")
	// DatabasePort := os.Getenv("DATABASE_PORT")

	// DatabaseURI := fmt.Sprintf(
	// 	"%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
	// 	DatabaseUsername,
	// 	DatabasePassword,
	// 	DatabaseHost,
	// 	DatabasePort,
	// 	DatabaseName,
	// )

	// db, err := gorm.Open("mysql", DatabaseURI)
	// if err != nil {
	// 	LOG.Println(err)
	// 	panic(err)
	// }
	// defer db.Close()

	// temporary SQLite for ease of development
	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&models.Mentor{})
	db.AutoMigrate(&models.Student{})
	db.AutoMigrate(&models.Project{})
}

// GetDB Get Database
func GetDB() *gorm.DB {
	// DatabaseUsername := os.Getenv("DATABASE_USERNAME")
	// DatabasePassword := os.Getenv("DATABASE_PASSWORD")
	// DatabaseName := os.Getenv("DATABASE_NAME")
	// DatabaseHost := os.Getenv("DATABASE_HOST")
	// DatabasePort := os.Getenv("DATABASE_PORT")

	// DatabaseURI := fmt.Sprintf(
	// 	"%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
	// 	DatabaseUsername,
	// 	DatabasePassword,
	// 	DatabaseHost,
	// 	DatabasePort,
	// 	DatabaseName,
	// )

	// db, err := gorm.Open("mysql", DatabaseURI)
	// if err != nil {
	// 	LOG.Println(err)
	// 	panic(err)
	// }

	// return db

	// temporary SQLite for ease of development
	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		LOG.Println(err)
		panic(err)
	}

	return db

}
