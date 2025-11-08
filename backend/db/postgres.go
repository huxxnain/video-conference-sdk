package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"video-conference-sdk/backend/models"
	"video-conference-sdk/backend/config"
)

var DB *gorm.DB

func SetupDatabase() {
	var err error
	DB, err = gorm.Open(postgres.Open(config.PostgresDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}
	err = DB.AutoMigrate(&models.Organization{}, &models.User{}, &models.Room{}, &models.QueueEntry{})
	if err != nil {
		log.Fatal("Migration error: ", err)
	}
}