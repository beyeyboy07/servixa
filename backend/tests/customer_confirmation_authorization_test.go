package tests

import (
	"strings"
	"testing"

	"servixa/internal/config"
	"servixa/internal/database"
	"servixa/internal/model"
	"servixa/internal/repository"
	"servixa/internal/service"
)

func TestCustomerConfirmationAuthorization(t *testing.T) {
	cfg := config.Load()
	db := database.Connect(cfg)

	if err := db.AutoMigrate(
		&model.CustomerConfirmation{},
	); err != nil {
		t.Fatalf("failed to migrate customer confirmation: %v", err)
	}

	// Create the customer who owns the work order.
	owner := &model.Customer{
		UserID:   999981,
		Name:     "Confirmation Owner",
		Email:    "confirmation-owner@servixa.local",
		Phone:    "0800000091",
		Address:  "Owner Test Address",
		IsActive: true,
	}
	if err := db.Create(owner).Error; err != nil {
		t.Fatalf("failed to create owner: %v", err)
	}
	defer db.Delete(&model.Customer{}, owner.ID)

	// Create another customer who must not access the work order.
	otherCustomer := &model.Customer{
		UserID:   999982,
		Name:     "Other Confirmation Customer",
		Email:    "other-confirmation@servixa.local",
		Phone:    "0800000092",
		Address:  "Other Test Address",
		IsActive: true,
	}
	if err := db.Create(otherCustomer).Error; err != nil {
		t.Fatalf("failed to create other customer: %v", err)
	}
	defer db.Delete(&model.Customer{}, otherCustomer.ID)

	workOrder := &model.WorkOrder{
		ServiceRequestID: 999981,
		CustomerID:       owner.ID,
		Title:            "Confirmation Authorization Test",
		Description:      "Authorization test",
		ServiceAddress:   owner.Address,
		Status:           model.WorkOrderCompleted,
	}
	if err := db.Create(workOrder).Error; err != nil {
		t.Fatalf("failed to create work order: %v", err)
	}
	defer db.Delete(&model.WorkOrder{}, workOrder.ID)

	completion := &model.WorkCompletion{
		WorkOrderID:     workOrder.ID,
		TechnicianID:    1,
		Summary:         "Authorization test completion",
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

	_, err := confirmationService.Confirm(
		workOrder.ID,
		otherCustomer.UserID,
		model.CustomerConfirmationConfirmed,
		"Unauthorized confirmation attempt",
	)

	if err == nil {
		t.Fatal("expected unauthorized customer confirmation to fail")
	}

	if !strings.Contains(err.Error(), "not authorized") {
		t.Fatalf("expected authorization error, got: %v", err)
	}

	var count int64
	if err := db.Model(&model.CustomerConfirmation{}).
		Where("work_order_id = ?", workOrder.ID).
		Count(&count).Error; err != nil {
		t.Fatalf("failed to count confirmations: %v", err)
	}

	if count != 0 {
		t.Fatalf("expected no confirmation to be created, got %d", count)
	}
}
