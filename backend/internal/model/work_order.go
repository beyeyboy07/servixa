package model

import "time"

type WorkOrderStatus string

const (
	WorkOrderPending            WorkOrderStatus = "PENDING"
	WorkOrderAssigned           WorkOrderStatus = "ASSIGNED"
	WorkOrderScheduled          WorkOrderStatus = "SCHEDULED"
	WorkOrderTechnicianAccepted WorkOrderStatus = "TECHNICIAN_ACCEPTED"
	WorkOrderInProgress         WorkOrderStatus = "IN_PROGRESS"
	WorkOrderCompleted          WorkOrderStatus = "COMPLETED"
	WorkOrderCustomerConfirmed  WorkOrderStatus = "CUSTOMER_CONFIRMED"
	WorkOrderClosed             WorkOrderStatus = "CLOSED"
	WorkOrderCancelled          WorkOrderStatus = "CANCELLED"
)

type WorkOrder struct {
	ID uint `gorm:"primaryKey"`

	ServiceRequestID uint `gorm:"not null;uniqueIndex"`
	CustomerID       uint `gorm:"not null;index"`

	Title       string `gorm:"type:varchar(150);not null"`
	Description string `gorm:"type:text;not null"`

	ServiceAddress string `gorm:"type:text;not null"`

	Status WorkOrderStatus `gorm:"type:varchar(30);not null;default:'PENDING'"`

	ScheduledStart *time.Time
	ScheduledEnd   *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
