package model

import "time"

type NotificationType string

const (
	NotificationWorkOrderCreated   NotificationType = "WORK_ORDER_CREATED"
	NotificationWorkOrderAssigned  NotificationType = "WORK_ORDER_ASSIGNED"
	NotificationWorkOrderScheduled NotificationType = "WORK_ORDER_SCHEDULED"
	NotificationWorkStarted        NotificationType = "WORK_STARTED"
	NotificationWorkCompleted      NotificationType = "WORK_COMPLETED"
	NotificationCustomerConfirmed  NotificationType = "CUSTOMER_CONFIRMED"
	NotificationWorkOrderCancelled NotificationType = "WORK_ORDER_CANCELLED"
)

type Notification struct {
	ID        uint             `gorm:"primaryKey"`
	UserID    uint             `gorm:"not null;index"`
	Type      NotificationType `gorm:"type:varchar(50);not null"`
	Title     string           `gorm:"type:varchar(150);not null"`
	Message   string           `gorm:"type:text;not null"`
	Reference string           `gorm:"type:varchar(100);index"`
	IsRead    bool             `gorm:"not null;default:false;index"`
	ReadAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
