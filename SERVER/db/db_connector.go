package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"server/models"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := ""
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	log.Println("✅ Connected to PostgreSQL")

	// Auto-migrate tables
	db.AutoMigrate(&models.Ping{})
}