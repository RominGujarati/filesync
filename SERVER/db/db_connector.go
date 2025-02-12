package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	DB_USER     = "your_user"
	DB_PASSWORD = "your_password"
	DB_NAME     = "your_database"
	DB_HOST     = "localhost"
	DB_PORT     = "5432"
)

// Open a database connection
func connectDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("✅ Connected to PostgreSQL")
	return db, nil
}

// Save Ping Notification
func savePing(db *sql.DB, status, details string) {
	_, err := db.Exec("INSERT INTO pings (status, details) VALUES ($1, $2)", status, details)
	if err != nil {
		log.Printf("❌ Failed to save ping: %v\n", err)
	} else {
		fmt.Println("📌 Ping saved successfully")
	}
}
