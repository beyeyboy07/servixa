package model

import "time"

type Customer struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"not null;uniqueIndex"`

	Name  string `gorm:"type:varchar(100);not null"`
	Email string `gorm:"type:varchar(150);not null"`
	Phone string `gorm:"type:varchar(30);not null"`

	Address string `gorm:"type:text;not null"`

	IsActive bool `gorm:"not null;default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
