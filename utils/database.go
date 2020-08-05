package utils

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"kwoc20-backend/models"
)

func InitialMigration() {
	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&models.Mentor{})
	db.AutoMigrate(&models.Project{})
}

func GetDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		LOG.Println(err)
		panic(err)
	}

	return db
}
	
