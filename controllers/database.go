package controllers

import "gorm.io/gorm"

type DBHandler struct {
	db *gorm.DB
}

func NewDBHandler(db *gorm.DB) *DBHandler {
	return &DBHandler{
		db: db,
	}
}
