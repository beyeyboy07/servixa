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

func TestWorkOrderVehicleConcurrency(t *testing.T) {
	cfg := config.Load()
	db := database.Connect(cfg)

	customer := &model.Customer{
		UserID:   1,
		Name:     "Vehicle Concurrency Customer",
		Email:    "vehicle-concurrency@servixa.local",
		Phone:    "0800000003",
		Address:  "Vehicle Test Address",
		IsActive: true,
	}

	if err := db.Create(customer).Error; err != nil {
		t.Fatalf("failed to create test customer: %v", err)
	}

	defer db.Delete(&model.Customer{}, customer.ID)

	workOrder := &model.WorkOrder{
		ServiceRequestID: 999998,
		CustomerID:       customer.ID,
		Title:            "Vehicle Concurrency Test",
		Description:      "Vehicle concurrency test",
		ServiceAddress:   "Vehicle Test Address",
		Status:           model.WorkOrderPending,
	}

	if err := db.Create(workOrder).Error; err != nil {
		t.Fatalf("failed to create test work order: %v", err)
	}

	defer db.Delete(&model.WorkOrder{}, workOrder.ID)

	vehicle := &model.Vehicle{
		PlateNumber: "B 9999 VCT",
		Brand:       "Test Brand",
		Model:       "Test Model",
		Year:        2025,
		Status:      "AVAILABLE",
		IsActive:    true,
	}

	if err := db.Create(vehicle).Error; err != nil {
		t.Fatalf("failed to create test vehicle: %v", err)
	}

	defer db.Delete(&model.Vehicle{}, vehicle.ID)

	workOrderVehicleRepository :=
		repository.NewWorkOrderVehicleRepository(db)

	workOrderRepository :=
		repository.NewWorkOrderRepository(db)

	vehicleRepository :=
		repository.NewVehicleRepository(db)

	workOrderVehicleService :=
		service.NewWorkOrderVehicleService(
			db,
			workOrderVehicleRepository,
			workOrderRepository,
			vehicleRepository,
		)

	var wg sync.WaitGroup
	results := make(chan error, 2)

	wg.Add(2)

	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()

			_, err := workOrderVehicleService.Assign(
				workOrder.ID,
				vehicle.ID,
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
		Model(&model.WorkOrderVehicle{}).
		Where("work_order_id = ?", workOrder.ID).
		Count(&assignmentCount).
		Error; err != nil {
		t.Fatalf(
			"failed to count vehicle assignments: %v",
			err,
		)
	}

	if assignmentCount != 1 {
		t.Fatalf(
			"expected exactly 1 vehicle assignment, got %d",
			assignmentCount,
		)
	}
}
