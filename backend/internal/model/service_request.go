package model

import "time"

type ServiceRequestStatus string

const (
	ServiceRequestPending  ServiceRequestStatus = "PENDING"
	ServiceRequestReviewed ServiceRequestStatus = "REVIEWED"
	ServiceRequestRejected ServiceRequestStatus = "REJECTED"
)

type ServiceRequest struct {
	ID uint `gorm:"primaryKey"`

	CustomerID uint `gorm:"not null;index"`

	Title       string `gorm:"type:varchar(150);not null"`
	Description string `gorm:"type:text;not null"`

	ServiceAddress string `gorm:"type:text;not null"`

	PreferredDate *time.Time

	Status ServiceRequestStatus `gorm:"type:varchar(20);not null;default:'PENDING'"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
