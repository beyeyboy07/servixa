package model

import "time"

type Technician struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"not null;uniqueIndex"`

	Name  string `gorm:"type:varchar(100);not null"`
	Email string `gorm:"type:varchar(150);not null"`
	Phone string `gorm:"type:varchar(30);not null"`

	Specialization string `gorm:"type:varchar(150);not null"`

	IsAvailable bool `gorm:"not null;default:true"`
	IsActive    bool `gorm:"not null;default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
