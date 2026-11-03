package tests

import (
	"testing"
	"time"

	"servixa/internal/config"
	"servixa/internal/database"
	"servixa/internal/model"
	"servixa/internal/repository"
	"servixa/internal/service"
)

func TestWorkOrderScheduleCreatesNotification(t *testing.T) {
	cfg := config.Load()
	db := database.Connect(cfg)

	// Siapkan teknisi.
	technician := &model.Technician{
		UserID:   1,
		Name:     "Schedule Test Technician",
		Email:    "schedule-test@servixa.local",
		Phone:    "081000000000",
		IsActive: true,
	}

	if err := db.Create(technician).Error; err != nil {
		t.Fatalf("failed to create technician: %v", err)
	}
	defer db.Delete(&model.Technician{}, technician.ID)

	// Siapkan kendaraan.
	vehicle := &model.Vehicle{
		PlateNumber: "TEST-SCHEDULE-001",
		Brand:       "Test Brand",
		Model:       "Test Model",
		IsActive:    true,
	}

	if err := db.Create(vehicle).Error; err != nil {
		t.Fatalf("failed to create vehicle: %v", err)
	}
	defer db.Delete(&model.Vehicle{}, vehicle.ID)

	// Siapkan customer dan service request.
	customer := &model.Customer{
		UserID:   1,
		Name:     "Schedule Test Customer",
		Email:    "schedule-customer@servixa.local",
		Phone:    "082000000000",
		Address:  "Test Address",
		IsActive: true,
	}

	if err := db.Create(customer).Error; err != nil {
		t.Fatalf("failed to create customer: %v", err)
	}
	defer db.Delete(&model.Customer{}, customer.ID)

	serviceRequest := &model.ServiceRequest{
		CustomerID:     customer.ID,
		Title:          "Schedule Test Work Order",
		Description:    "Test work order scheduling",
		ServiceAddress: "Test Service Address",
		Status:         model.ServiceRequestReviewed,
	}

	if err := db.Create(serviceRequest).Error; err != nil {
		t.Fatalf("failed to create service request: %v", err)
	}
	defer db.Delete(&model.ServiceRequest{}, serviceRequest.ID)

	// Siapkan work order dalam status ASSIGNED.
	workOrder := &model.WorkOrder{
		ServiceRequestID: serviceRequest.ID,
		Title:            serviceRequest.Title,
		Status:           model.WorkOrderAssigned,
	}

	if err := db.Create(workOrder).Error; err != nil {
		t.Fatalf("failed to create work order: %v", err)
	}
	defer db.Delete(&model.WorkOrder{}, workOrder.ID)

	// Siapkan assignment teknisi.
	assignment := &model.WorkOrderAssignment{
		WorkOrderID:  workOrder.ID,
		TechnicianID: technician.ID,
	}

	if err := db.Create(assignment).Error; err != nil {
		t.Fatalf("failed to create work order assignment: %v", err)
	}
	defer db.Delete(&model.WorkOrderAssignment{}, assignment.ID)

	// Siapkan assignment kendaraan.
	vehicleAssignment := &model.WorkOrderVehicle{
		WorkOrderID: workOrder.ID,
		VehicleID:   vehicle.ID,
	}

	if err := db.Create(vehicleAssignment).Error; err != nil {
		t.Fatalf("failed to create work order vehicle: %v", err)
	}
	defer db.Delete(&model.WorkOrderVehicle{}, vehicleAssignment.ID)

	// Inisialisasi repository dan service.
	scheduleRepository := repository.NewWorkOrderScheduleRepository(db)
	workOrderRepository := repository.NewWorkOrderRepository(db)
	assignmentRepository := repository.NewWorkOrderAssignmentRepository(db)
	vehicleRepository := repository.NewWorkOrderVehicleRepository(db)
	technicianRepository := repository.NewTechnicianRepository(db)
	notificationRepository := repository.NewNotificationRepository(db)

	scheduleService := service.NewWorkOrderScheduleService(
		db,
		scheduleRepository,
		workOrderRepository,
		assignmentRepository,
		vehicleRepository,
		technicianRepository,
		notificationRepository,
	)

	start := time.Now().Add(48 * time.Hour).Truncate(time.Minute)
	end := start.Add(2 * time.Hour)

	// Jalankan penjadwalan.
	schedule, err := scheduleService.Schedule(workOrder.ID, start, end)
	if err != nil {
		t.Fatalf("expected schedule to succeed, got error: %v", err)
	}

	if schedule == nil {
		t.Fatal("expected schedule, got nil")
	}

	// Pastikan status work order berubah.
	var updatedWorkOrder model.WorkOrder
	if err := db.First(&updatedWorkOrder, workOrder.ID).Error; err != nil {
		t.Fatalf("failed to load work order: %v", err)
	}

	if updatedWorkOrder.Status != model.WorkOrderScheduled {
		t.Fatalf(
			"expected work order status %v, got %v",
			model.WorkOrderScheduled,
			updatedWorkOrder.Status,
		)
	}

	// Pastikan notifikasi teknisi tercatat.
	var notification model.Notification
	if err := db.Where(
		"user_id = ? AND type = ? AND reference = ?",
		technician.UserID,
		model.NotificationWorkOrderScheduled,
		"WORK_ORDER:"+formatUint(workOrder.ID),
	).First(&notification).Error; err != nil {
		t.Fatalf("expected scheduling notification, got error: %v", err)
	}

	if notification.IsRead {
		t.Fatal("expected notification to be unread")
	}

	if notification.Title != "Work Order Scheduled" {
		t.Errorf("unexpected notification title: %s", notification.Title)
	}

	// Bersihkan data notifikasi dan jadwal yang dibuat test.
	defer db.Delete(&model.Notification{}, notification.ID)
	defer db.Delete(&model.WorkOrderSchedule{}, schedule.ID)
}

func TestWorkOrderScheduleRejectsInvalidTime(t *testing.T) {
	cfg := config.Load()
	db := database.Connect(cfg)

	scheduleService := service.NewWorkOrderScheduleService(
		db,
		repository.NewWorkOrderScheduleRepository(db),
		repository.NewWorkOrderRepository(db),
		repository.NewWorkOrderAssignmentRepository(db),
		repository.NewWorkOrderVehicleRepository(db),
		repository.NewTechnicianRepository(db),
		repository.NewNotificationRepository(db),
	)

	start := time.Now().Add(48 * time.Hour)
	end := start

	_, err := scheduleService.Schedule(1, start, end)
	if err == nil {
		t.Fatal("expected invalid schedule time to return an error")
	}
}
