package tests

import (
	"sync"
	"testing"

	"servixa/internal/config"
	"servixa/internal/database"
	"servixa/internal/model"
	"servixa/internal/repository"
	"servixa/internal/service"
)

func TestCustomerConfirmationConcurrency(t *testing.T) {
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
		UserID:   999983,
		Name:     "Confirmation Concurrency Customer",
		Email:    "confirmation-concurrency@servixa.local",
		Phone:    "0800000093",
		Address:  "Confirmation Test Address",
		IsActive: true,
	}

	if err := db.Create(customer).Error; err != nil {
		t.Fatalf("failed to create test customer: %v", err)
	}
	defer db.Delete(&model.Customer{}, customer.ID)

	technician := &model.Technician{
		UserID:         999986,
		Name:           "Confirmation Concurrency Technician",
		Email:          "confirmation-concurrency-tech@servixa.local",
		Phone:          "0800000096",
		Specialization: "General Maintenance",
		IsAvailable:    true,
		IsActive:       true,
	}

	if err := db.Create(technician).Error; err != nil {
		t.Fatalf("failed to create test technician: %v", err)
	}
	defer db.Delete(&model.Technician{}, technician.ID)

	workOrder := &model.WorkOrder{
		ServiceRequestID: 999983,
		CustomerID:       customer.ID,
		Title:            "Confirmation Concurrency Test",
		Description:      "Confirmation concurrency test",
		ServiceAddress:   customer.Address,
		Status:           model.WorkOrderCompleted,
	}

	if err := db.Create(workOrder).Error; err != nil {
		t.Fatalf("failed to create test work order: %v", err)
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
		Summary:         "Concurrency test completion",
		WorkPerformed:   "Required service work completed.",
		Recommendations: "Perform regular maintenance.",
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

	var wg sync.WaitGroup
	results := make(chan error, 2)

	wg.Add(2)

	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()

			_, err := confirmationService.Confirm(
				workOrder.ID,
				customer.UserID,
				model.CustomerConfirmationConfirmed,
				"Concurrency confirmation test",
			)

			results <- err
		}()
	}

	wg.Wait()
	close(results)

	successCount := 0
	errorCount := 0

	for err := range results {
		if err == nil {
			successCount++
		} else {
			errorCount++
		}
	}

	if successCount != 1 {
		t.Fatalf(
			"expected exactly 1 successful confirmation, got %d",
			successCount,
		)
	}

	if errorCount != 1 {
		t.Fatalf(
			"expected exactly 1 failed confirmation, got %d",
			errorCount,
		)
	}

	var confirmationCount int64

	if err := db.
		Model(&model.CustomerConfirmation{}).
		Where("work_order_id = ?", workOrder.ID).
		Count(&confirmationCount).
		Error; err != nil {
		t.Fatalf("failed to count confirmations: %v", err)
	}

	if confirmationCount != 1 {
		t.Fatalf(
			"expected exactly 1 customer confirmation, got %d",
			confirmationCount,
		)
	}

	var notificationCount int64

	if err := db.
		Model(&model.Notification{}).
		Where(
			"user_id = ? AND type = ? AND reference = ?",
			technician.UserID,
			model.NotificationCustomerConfirmed,
			"WORK_ORDER:"+formatUint(workOrder.ID),
		).
		Count(&notificationCount).
		Error; err != nil {
		t.Fatalf("failed to count notifications: %v", err)
	}

	if notificationCount != 1 {
		t.Fatalf(
			"expected exactly 1 customer confirmation notification, got %d",
			notificationCount,
		)
	}

	defer db.
		Where("user_id = ?", technician.UserID).
		Where("type = ?", model.NotificationCustomerConfirmed).
		Where("reference = ?", "WORK_ORDER:"+formatUint(workOrder.ID)).
		Delete(&model.Notification{})
}
