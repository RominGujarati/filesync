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

	// Run migrations
	migrateDB(db)
}

func migrateDB(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Ping{},
		&models.Job{},
		&models.JobStep{},
		&models.JobResult{},
		&models.User{},
		&models.YAMLUpload{},
		&models.AuditLog{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("✅ Database migrated successfully")
}