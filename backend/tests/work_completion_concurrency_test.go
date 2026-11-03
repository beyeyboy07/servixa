package tests

import (
	"strconv"
	"sync"
	"testing"

	"servixa/internal/config"
	"servixa/internal/database"
	"servixa/internal/model"
	"servixa/internal/repository"
	"servixa/internal/service"
)

func TestWorkCompletionConcurrency(t *testing.T) {
	cfg := config.Load()
	db := database.Connect(cfg)

	if err := db.AutoMigrate(
		&model.Customer{},
		&model.ServiceRequest{},
		&model.WorkOrder{},
		&model.Technician{},
		&model.WorkOrderAssignment{},
		&model.WorkCompletion{},
		&model.ServiceVisit{},
		&model.Notification{},
	); err != nil {
		t.Fatalf("failed to migrate test tables: %v", err)
	}

	customer := &model.Customer{
		UserID:   1,
		Name:     "Completion Concurrency Customer",
		Email:    "completion-concurrency@servixa.local",
		Phone:    "0800000003",
		Address:  "Completion Test Address",
		IsActive: true,
	}

	if err := db.Create(customer).Error; err != nil {
		t.Fatalf("failed to create test customer: %v", err)
	}
	defer db.Delete(&model.Customer{}, customer.ID)

	serviceRequest := &model.ServiceRequest{
		CustomerID:     customer.ID,
		Title:          "Completion Concurrency Request",
		Description:    "Completion concurrency test",
		ServiceAddress: "Completion Test Address",
		Status:         model.ServiceRequestReviewed,
	}

	if err := db.Create(serviceRequest).Error; err != nil {
		t.Fatalf("failed to create service request: %v", err)
	}
	defer db.Delete(&model.ServiceRequest{}, serviceRequest.ID)

	workOrder := &model.WorkOrder{
		ServiceRequestID: serviceRequest.ID,
		CustomerID:       customer.ID,
		Title:            "Completion Concurrency Test",
		Description:      "Completion concurrency test",
		ServiceAddress:   "Completion Test Address",
		Status:           model.WorkOrderInProgress,
	}

	if err := db.Create(workOrder).Error; err != nil {
		t.Fatalf("failed to create test work order: %v", err)
	}
	defer db.Delete(&model.WorkOrder{}, workOrder.ID)

	technician := &model.Technician{
		UserID:         999998,
		Name:           "Completion Concurrency Technician",
		Email:          "completion-tech@servixa.local",
		Phone:          "0800000004",
		Specialization: "General Service",
		IsAvailable:    true,
		IsActive:       true,
	}

	if err := db.Create(technician).Error; err != nil {
		t.Fatalf("failed to create test technician: %v", err)
	}
	defer db.Delete(&model.Technician{}, technician.ID)

	assignment := &model.WorkOrderAssignment{
		WorkOrderID:  workOrder.ID,
		TechnicianID: technician.ID,
	}

	if err := db.Create(assignment).Error; err != nil {
		t.Fatalf("failed to create test assignment: %v", err)
	}
	defer db.Delete(&model.WorkOrderAssignment{}, assignment.ID)

	visit := &model.ServiceVisit{
		WorkOrderID:  workOrder.ID,
		TechnicianID: technician.ID,
		Status:       model.ServiceVisitCompleted,
	}

	if err := db.Create(visit).Error; err != nil {
		t.Fatalf("failed to create completed service visit: %v", err)
	}
	defer db.Delete(&model.ServiceVisit{}, visit.ID)

	workCompletionService := service.NewWorkCompletionService(
		db,
		repository.NewWorkCompletionRepository(db),
		repository.NewWorkOrderRepository(db),
		repository.NewWorkOrderAssignmentRepository(db),
		repository.NewTechnicianRepository(db),
		repository.NewServiceVisitRepository(db),
		repository.NewServiceRequestRepository(db),
		repository.NewCustomerRepository(db),
		repository.NewNotificationRepository(db),
	)

	var wg sync.WaitGroup
	results := make(chan error, 2)

	wg.Add(2)

	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()

			_, err := workCompletionService.Create(
				workOrder.ID,
				technician.UserID,
				"Completion concurrency test",
				"Technician performed the required service work.",
				"Perform regular maintenance.",
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
			"expected exactly 1 successful completion, got %d",
			successCount,
		)
	}

	if errorCount != 1 {
		t.Fatalf(
			"expected exactly 1 failed completion, got %d",
			errorCount,
		)
	}

	var completionCount int64

	if err := db.
		Model(&model.WorkCompletion{}).
		Where("work_order_id = ?", workOrder.ID).
		Count(&completionCount).
		Error; err != nil {
		t.Fatalf("failed to count work completions: %v", err)
	}

	if completionCount != 1 {
		t.Fatalf(
			"expected exactly 1 work completion, got %d",
			completionCount,
		)
	}

	var updatedWorkOrder model.WorkOrder

	if err := db.First(&updatedWorkOrder, workOrder.ID).Error; err != nil {
		t.Fatalf("failed to retrieve work order: %v", err)
	}

	if updatedWorkOrder.Status != model.WorkOrderCompleted {
		t.Fatalf(
			"expected work order status %q, got %q",
			model.WorkOrderCompleted,
			updatedWorkOrder.Status,
		)
	}

	var notificationCount int64

	if err := db.
		Model(&model.Notification{}).
		Where(
			"user_id = ? AND reference = ? AND type = ?",
			customer.UserID,
			"WORK_ORDER:"+uintToString(workOrder.ID),
			model.NotificationWorkCompleted,
		).
		Count(&notificationCount).
		Error; err != nil {
		t.Fatalf("failed to count completion notifications: %v", err)
	}

	if notificationCount != 1 {
		t.Fatalf(
			"expected exactly 1 completion notification, got %d",
			notificationCount,
		)
	}
}

func uintToString(value uint) string {
	return strconv.FormatUint(uint64(value), 10)
}
