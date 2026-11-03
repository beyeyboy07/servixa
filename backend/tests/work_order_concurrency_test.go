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

func TestWorkOrderConcurrency(t *testing.T) {
	cfg := config.Load()
	db := database.Connect(cfg)

	customer := &model.Customer{
		UserID:   1,
		Name:     "Concurrency Test Customer",
		Email:    "concurrency-test@servixa.local",
		Phone:    "0800000000",
		Address:  "Test Address",
		IsActive: true,
	}

	if err := db.Create(customer).Error; err != nil {
		t.Fatalf("failed to create test customer: %v", err)
	}

	defer db.Delete(&model.Customer{}, customer.ID)

	serviceRequest := &model.ServiceRequest{
		CustomerID:     customer.ID,
		Title:          "Concurrency Test",
		Description:    "Concurrency test service request",
		ServiceAddress: "Test Service Address",
		Status:         model.ServiceRequestReviewed,
	}

	if err := db.Create(serviceRequest).Error; err != nil {
		t.Fatalf("failed to create test service request: %v", err)
	}

	defer db.Delete(&model.ServiceRequest{}, serviceRequest.ID)

	workOrderRepository := repository.NewWorkOrderRepository(db)
	serviceRequestRepository := repository.NewServiceRequestRepository(db)
	customerRepository := repository.NewCustomerRepository(db)
	notificationRepository := repository.NewNotificationRepository(db)

	workOrderAssignmentRepository :=
		repository.NewWorkOrderAssignmentRepository(db)

	technicianRepository :=
		repository.NewTechnicianRepository(db)

	workOrderService := service.NewWorkOrderService(
		db,
		workOrderRepository,
		serviceRequestRepository,
		customerRepository,
		notificationRepository,
		workOrderAssignmentRepository,
		technicianRepository,
	)

	var wg sync.WaitGroup
	results := make(chan error, 2)

	wg.Add(2)

	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()

			_, err := workOrderService.Create(serviceRequest.ID)
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
			"expected exactly 1 successful creation, got %d",
			successCount,
		)
	}

	if errorCount != 1 {
		t.Fatalf(
			"expected exactly 1 failed creation, got %d",
			errorCount,
		)
	}

	var workOrder model.WorkOrder
	if err := db.
		Where("service_request_id = ?", serviceRequest.ID).
		First(&workOrder).
		Error; err != nil {
		t.Fatalf("failed to find created work order: %v", err)
	}

	defer db.Where(
		"reference = ?",
		"WORK_ORDER:"+formatUint(workOrder.ID),
	).Delete(&model.Notification{})

	defer db.Delete(&model.WorkOrder{}, workOrder.ID)

	var workOrderCount int64

	if err := db.
		Model(&model.WorkOrder{}).
		Where("service_request_id = ?", serviceRequest.ID).
		Count(&workOrderCount).
		Error; err != nil {
		t.Fatalf("failed to count work orders: %v", err)
	}

	if workOrderCount != 1 {
		t.Fatalf(
			"expected exactly 1 work order, got %d",
			workOrderCount,
		)
	}

	var notificationCount int64

	if err := db.
		Model(&model.Notification{}).
		Where(
			"reference = ?",
			"WORK_ORDER:"+formatUint(workOrder.ID),
		).
		Count(&notificationCount).
		Error; err != nil {
		t.Fatalf("failed to count notifications: %v", err)
	}

	if notificationCount != 1 {
		t.Fatalf(
			"expected exactly 1 notification, got %d",
			notificationCount,
		)
	}
}

func formatUint(value uint) string {
	return strconv.FormatUint(uint64(value), 10)
}
