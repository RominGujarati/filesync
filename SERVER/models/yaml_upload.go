package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Username     string    `gorm:"unique;not null"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type YAMLUpload struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null"`
	FileName   string    `gorm:"not null"`
	// Version    int       `gorm:"not null"`
	Checksum   string    `gorm:"not null"`
	UploadedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Content    string    `gorm:"not null"`
}
