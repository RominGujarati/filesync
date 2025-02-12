package models

import (
	"time"
	// "gorm.io/gorm"
)

type Ping struct {
	ID         uint           `gorm:"primaryKey"`
	ReceivedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	Status     string         `gorm:"not null"`
	Details    string
	IP         string
}
