package model

import "time"

type ServiceVisitStatus string

const (
	ServiceVisitScheduled  ServiceVisitStatus = "SCHEDULED"
	ServiceVisitInProgress ServiceVisitStatus = "IN_PROGRESS"
	ServiceVisitCompleted  ServiceVisitStatus = "COMPLETED"
)

type ServiceVisit struct {
	ID uint `gorm:"primaryKey"`

	WorkOrderID  uint `gorm:"not null;uniqueIndex"`
	TechnicianID uint `gorm:"not null;index"`

	Status ServiceVisitStatus `gorm:"type:varchar(30);not null;default:'SCHEDULED'"`

	StartedAt   *time.Time
	CompletedAt *time.Time

	Notes string `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
