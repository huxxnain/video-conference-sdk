package db

import (
	"log"
	"video-conference-sdk/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func InitPostgres() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://user:password@postgres:5432/vconfdb?sslmode=disable"
	}
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB = database

	// Auto migrate models
	if err := DB.AutoMigrate(
		&models.Organization{},
		&models.User{},
		&models.Room{},
		&models.QueueEntry{},
	); err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	
	log.Println("Database connected and migrated successfully")
}