package models

import (
	"time"
)

type AuditLog struct {
	ID             uint      `gorm:"primaryKey"`
	UserID         uint      `gorm:"not null"`
	ActionType     string    `gorm:"not null"`
	TargetType     string    `gorm:"not null"`
	TargetID       uint
	ActionTimestamp time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Details        string
}
