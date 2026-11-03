package model

import "time"

type UserRole string

const (
	RoleAdmin      UserRole = "ADMIN"
	RoleDispatcher UserRole = "DISPATCHER"
	RoleTechnician UserRole = "TECHNICIAN"
	RoleCustomer   UserRole = "CUSTOMER"
)

type User struct {
	ID        uint     `gorm:"primaryKey"`
	Name      string   `gorm:"type:varchar(100);not null"`
	Email     string   `gorm:"type:varchar(150);not null;uniqueIndex"`
	Password  string   `gorm:"type:varchar(255);not null"`
	Role      UserRole `gorm:"type:varchar(20);not null"`
	IsActive  bool     `gorm:"not null;default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
