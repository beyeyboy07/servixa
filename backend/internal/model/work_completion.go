package model

import "time"

type WorkCompletion struct {
	ID              uint   `gorm:"primaryKey"`
	WorkOrderID     uint   `gorm:"not null;uniqueIndex"`
	TechnicianID    uint   `gorm:"not null;index"`
	Summary         string `gorm:"type:text;not null"`
	WorkPerformed   string `gorm:"type:text;not null"`
	Recommendations string `gorm:"type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
