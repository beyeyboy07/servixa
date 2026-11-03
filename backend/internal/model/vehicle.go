package model

import "time"

type Vehicle struct {
	ID uint `gorm:"primaryKey"`

	PlateNumber string `gorm:"type:varchar(20);not null;uniqueIndex"`
	Brand       string `gorm:"type:varchar(50);not null"`
	Model       string `gorm:"type:varchar(50);not null"`
	Year        int    `gorm:"not null"`

	Status string `gorm:"type:varchar(20);not null;default:'AVAILABLE'"`

	IsActive bool `gorm:"not null;default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateVehicleRequest struct {
	PlateNumber string `json:"plate_number"`
	Brand       string `json:"brand"`
	Model       string `json:"model"`
	Year        int    `json:"year"`
	Status      string `json:"status"`
}
