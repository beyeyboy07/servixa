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

func TestWorkOrderAssignmentConcurrency(t *testing.T) {
	cfg := config.Load()
	db := database.Connect(cfg)

	customer := &model.Customer{
		UserID:   1,
		Name:     "Assignment Concurrency Customer",
		Email:    "assignment-concurrency@servixa.local",
		Phone:    "0800000001",
		Address:  "Assignment Test Address",
		IsActive: true,
	}

	if err := db.Create(customer).Error; err != nil {
		t.Fatalf("failed to create test customer: %v", err)
	}

	defer db.Delete(&model.Customer{}, customer.ID)

	workOrder := &model.WorkOrder{
		ServiceRequestID: 999999,
		CustomerID:       customer.ID,
		Title:            "Assignment Concurrency Test",
		Description:      "Assignment concurrency test",
		ServiceAddress:   "Assignment Test Address",
		Status:           model.WorkOrderPending,
	}

	if err := db.Create(workOrder).Error; err != nil {
		t.Fatalf("failed to create test work order: %v", err)
	}

	defer db.Delete(&model.WorkOrder{}, workOrder.ID)

	technician := &model.Technician{
		UserID:         1,
		Name:           "Assignment Concurrency Technician",
		Email:          "assignment-tech@servixa.local",
		Phone:          "0800000002",
		Specialization: "General Service",
		IsAvailable:    true,
		IsActive:       true,
	}

	if err := db.Create(technician).Error; err != nil {
		t.Fatalf("failed to create test technician: %v", err)
	}

	defer db.Delete(&model.Technician{}, technician.ID)

	workOrderAssignmentRepository :=
		repository.NewWorkOrderAssignmentRepository(db)

	workOrderRepository :=
		repository.NewWorkOrderRepository(db)

	technicianRepository :=
		repository.NewTechnicianRepository(db)

	notificationRepository :=
		repository.NewNotificationRepository(db)

	workOrderAssignmentService :=
		service.NewWorkOrderAssignmentService(
			db,
			workOrderAssignmentRepository,
			workOrderRepository,
			technicianRepository,
			notificationRepository,
		)

	var wg sync.WaitGroup
	results := make(chan error, 2)

	wg.Add(2)

	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()

			_, err := workOrderAssignmentService.Assign(
				workOrder.ID,
				technician.ID,
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
			continue
		}

		errorCount++
	}

	if successCount != 1 {
		t.Fatalf(
			"expected exactly 1 successful assignment, got %d",
			successCount,
		)
	}

	if errorCount != 1 {
		t.Fatalf(
			"expected exactly 1 failed assignment, got %d",
			errorCount,
		)
	}

	var assignmentCount int64

	if err := db.
		Model(&model.WorkOrderAssignment{}).
		Where("work_order_id = ?", workOrder.ID).
		Count(&assignmentCount).
		Error; err != nil {
		t.Fatalf("failed to count assignments: %v", err)
	}

	if assignmentCount != 1 {
		t.Fatalf(
			"expected exactly 1 assignment, got %d",
			assignmentCount,
		)
	}

	var updatedWorkOrder model.WorkOrder

	if err := db.
		Where("id = ?", workOrder.ID).
		First(&updatedWorkOrder).
		Error; err != nil {
		t.Fatalf("failed to fetch work order: %v", err)
	}

	if updatedWorkOrder.Status != model.WorkOrderAssigned {
		t.Fatalf(
			"expected work order status %s, got %s",
			model.WorkOrderAssigned,
			updatedWorkOrder.Status,
		)
	}
}
