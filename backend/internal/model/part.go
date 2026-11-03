package model

import "time"

type Part struct {
	ID uint `gorm:"primaryKey"`

	Code string `gorm:"type:varchar(50);not null;uniqueIndex"`
	Name string `gorm:"type:varchar(150);not null"`

	Description string `gorm:"type:text"`

	Unit string `gorm:"type:varchar(20);not null"`

	Stock    int `gorm:"not null;default:0"`
	MinStock int `gorm:"not null;default:0"`

	UnitPrice float64 `gorm:"type:decimal(15,2);not null;default:0"`

	IsActive bool `gorm:"not null;default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
