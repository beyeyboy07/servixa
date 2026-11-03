package tests

import (
	"testing"

	"servixa/internal/config"
	"servixa/internal/database"
	"servixa/internal/model"
	"servixa/internal/repository"
	"servixa/internal/service"
)

func TestCustomerConfirmationNotification(t *testing.T) {
	cfg := config.Load()
	db := database.Connect(cfg)

	if err := db.AutoMigrate(
		&model.Customer{},
		&model.Technician{},
		&model.WorkOrder{},
		&model.WorkOrderAssignment{},
		&model.WorkCompletion{},
		&model.CustomerConfirmation{},
		&model.Notification{},
	); err != nil {
		t.Fatalf("failed to migrate tables: %v", err)
	}

	customer := &model.Customer{
		UserID:   999984,
		Name:     "Confirmation Notification Customer",
		Email:    "confirmation-notification@servixa.local",
		Phone:    "0800000094",
		Address:  "Confirmation Notification Address",
		IsActive: true,
	}
	if err := db.Create(customer).Error; err != nil {
		t.Fatalf("failed to create customer: %v", err)
	}
	defer db.Delete(&model.Customer{}, customer.ID)

	technician := &model.Technician{
		UserID:         999985,
		Name:           "Confirmation Notification Technician",
		Email:          "confirmation-technician@servixa.local",
		Phone:          "0800000095",
		Specialization: "General Maintenance",
		IsAvailable:    true,
		IsActive:       true,
	}
	if err := db.Create(technician).Error; err != nil {
		t.Fatalf("failed to create technician: %v", err)
	}
	defer db.Delete(&model.Technician{}, technician.ID)

	workOrder := &model.WorkOrder{
		ServiceRequestID: 999984,
		CustomerID:       customer.ID,
		Title:            "Customer Confirmation Notification Test",
		Description:      "Test customer confirmation notification",
		ServiceAddress:   customer.Address,
		Status:           model.WorkOrderCompleted,
	}
	if err := db.Create(workOrder).Error; err != nil {
		t.Fatalf("failed to create work order: %v", err)
	}
	defer db.Delete(&model.WorkOrder{}, workOrder.ID)

	assignment := &model.WorkOrderAssignment{
		WorkOrderID:  workOrder.ID,
		TechnicianID: technician.ID,
	}
	if err := db.Create(assignment).Error; err != nil {
		t.Fatalf("failed to create work order assignment: %v", err)
	}
	defer db.Delete(&model.WorkOrderAssignment{}, assignment.ID)

	completion := &model.WorkCompletion{
		WorkOrderID:     workOrder.ID,
		TechnicianID:    technician.ID,
		Summary:         "Notification test completion",
		WorkPerformed:   "Service work completed",
		Recommendations: "Regular maintenance",
	}
	if err := db.Create(completion).Error; err != nil {
		t.Fatalf("failed to create work completion: %v", err)
	}
	defer db.Delete(&model.WorkCompletion{}, completion.ID)

	confirmationRepository :=
		repository.NewCustomerConfirmationRepository(db)
	workOrderRepository :=
		repository.NewWorkOrderRepository(db)
	workCompletionRepository :=
		repository.NewWorkCompletionRepository(db)
	customerRepository :=
		repository.NewCustomerRepository(db)
	workOrderAssignmentRepository :=
		repository.NewWorkOrderAssignmentRepository(db)
	technicianRepository :=
		repository.NewTechnicianRepository(db)
	notificationRepository :=
		repository.NewNotificationRepository(db)

	confirmationService := service.NewCustomerConfirmationService(
		db,
		confirmationRepository,
		workOrderRepository,
		workCompletionRepository,
		customerRepository,
		workOrderAssignmentRepository,
		technicianRepository,
		notificationRepository,
	)

	confirmation, err := confirmationService.Confirm(
		workOrder.ID,
		customer.UserID,
		model.CustomerConfirmationConfirmed,
		"Work completed successfully",
	)
	if err != nil {
		t.Fatalf("failed to confirm work order: %v", err)
	}

	if confirmation.Status != model.CustomerConfirmationConfirmed {
		t.Fatalf(
			"expected confirmation status %v, got %v",
			model.CustomerConfirmationConfirmed,
			confirmation.Status,
		)
	}

	// The work order should be closed after confirmation.
	var updatedWorkOrder model.WorkOrder
	if err := db.First(&updatedWorkOrder, workOrder.ID).Error; err != nil {
		t.Fatalf("failed to reload work order: %v", err)
	}

	if updatedWorkOrder.Status != model.WorkOrderClosed {
		t.Fatalf(
			"expected work order status %v, got %v",
			model.WorkOrderClosed,
			updatedWorkOrder.Status,
		)
	}

	// The notification must belong to the assigned technician.
	var notifications []model.Notification
	if err := db.
		Where(
			"user_id = ? AND type = ? AND reference = ?",
			technician.UserID,
			model.NotificationCustomerConfirmed,
			"WORK_ORDER:"+formatUint(workOrder.ID),
		).
		Find(&notifications).Error; err != nil {
		t.Fatalf("failed to query notifications: %v", err)
	}

	if len(notifications) != 1 {
		t.Fatalf(
			"expected exactly 1 customer confirmation notification, got %d",
			len(notifications),
		)
	}

	notification := notifications[0]

	if notification.Title != "Customer Confirmed Work" {
		t.Errorf(
			"unexpected notification title: %s",
			notification.Title,
		)
	}

	if notification.IsRead {
		t.Error("expected notification to be unread")
	}

	defer db.Delete(&model.Notification{}, notification.ID)
}
