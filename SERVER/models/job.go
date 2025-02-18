package models

import (
	"time"
)

type Job struct {
	ID          uint      `gorm:"primaryKey"`
	JobName     string    `gorm:"not null"`
	Description string
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	StartedAt   *time.Time
	CompletedAt *time.Time
	Status      string    `gorm:"not null"`
	Result      string
	Steps       []JobStep `gorm:"foreignKey:JobID"`
}

type JobStep struct {
	ID          uint      `gorm:"primaryKey"`
	JobID       uint      `gorm:"not null"`
	StepName    string    `gorm:"not null"`
	Description string
	StartedAt   *time.Time
	CompletedAt *time.Time
	Status      string    `gorm:"not null"`
	Result      string
}

type JobResult struct {
	ID           uint      `gorm:"primaryKey"`
	JobID        uint      `gorm:"not null"`
	StepID       *uint
	Status       string    `gorm:"not null"`
	Output       string
	ErrorDetails string
	RecordedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
