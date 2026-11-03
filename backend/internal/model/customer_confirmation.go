package model

import "time"

type CustomerConfirmationStatus string

const (
	CustomerConfirmationConfirmed CustomerConfirmationStatus = "CONFIRMED"
	CustomerConfirmationRejected  CustomerConfirmationStatus = "REJECTED"
)

type CustomerConfirmation struct {
	ID          uint                       `gorm:"primaryKey"`
	WorkOrderID uint                       `gorm:"not null;uniqueIndex"`
	CustomerID  uint                       `gorm:"not null;index"`
	Status      CustomerConfirmationStatus `gorm:"type:varchar(30);not null"`
	Notes       string                     `gorm:"type:text"`
	ConfirmedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
