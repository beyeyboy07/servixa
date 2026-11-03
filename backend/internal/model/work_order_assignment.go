package model

import "time"

type WorkOrderAssignment struct {
	ID uint `gorm:"primaryKey"`

	WorkOrderID  uint `gorm:"not null;uniqueIndex"`
	TechnicianID uint `gorm:"not null;index"`

	AssignedAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
