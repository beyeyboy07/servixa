package model

import "time"

type WorkOrderSchedule struct {
	ID uint `gorm:"primaryKey"`

	WorkOrderID    uint      `gorm:"not null;uniqueIndex"`
	ScheduledStart time.Time `gorm:"not null"`
	ScheduledEnd   time.Time `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
