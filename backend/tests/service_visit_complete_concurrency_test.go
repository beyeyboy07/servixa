package tests

import (
	"fmt"
	"testing"
	"time"

	"servixa/internal/config"
	"servixa/internal/database"
	"servixa/internal/model"
	"servixa/internal/repository"
	"servixa/internal/service"
)

func TestServiceVisitCompleteCreatesNotification(t *testing.T) {
	cfg := config.Load()
	db := database.Connect(cfg)

	customer := &model.Customer{
		UserID:   1,
		Name:     "Complete Test Customer",
		Email:    "complete-customer@servixa.local",
		Phone:    "085000000000",
		Address:  "Test Address",
		IsActive: true,
	}
	if err := db.Create(customer).Error; err != nil {
		t.Fatalf("failed to create customer: %v", err)
	}
	defer db.Delete(&model.Customer{}, customer.ID)

	serviceRequest := &model.ServiceRequest{
		CustomerID:     customer.ID,
		Title:          "Complete Test Work Order",
		Description:    "Testing service visit completion",
		ServiceAddress: "Test Service Address",
		Status:         model.ServiceRequestReviewed,
	}
	if err := db.Create(serviceRequest).Error; err != nil {
		t.Fatalf("failed to create service request: %v", err)
	}
	defer db.Delete(&model.ServiceRequest{}, serviceRequest.ID)

	workOrder := &model.WorkOrder{
		ServiceRequestID: serviceRequest.ID,
		Title:            serviceRequest.Title,
		Status:           model.WorkOrderInProgress,
	}
	if err := db.Create(workOrder).Error; err != nil {
		t.Fatalf("failed to create work order: %v", err)
	}
	defer db.Delete(&model.WorkOrder{}, workOrder.ID)

	technician := &model.Technician{
		UserID:   1,
		Name:     "Complete Test Technician",
		Email:    "complete-tech@servixa.local",
		Phone:    "086000000000",
		IsActive: true,
	}
	if err := db.Create(technician).Error; err != nil {
		t.Fatalf("failed to create technician: %v", err)
	}
	defer db.Delete(&model.Technician{}, technician.ID)

	startedAt := time.Now().Add(-time.Hour)
	visit := &model.ServiceVisit{
		WorkOrderID:  workOrder.ID,
		TechnicianID: technician.ID,
		Status:       model.ServiceVisitInProgress,
		StartedAt:    &startedAt,
	}
	if err := db.Create(visit).Error; err != nil {
		t.Fatalf("failed to create service visit: %v", err)
	}
	defer db.Delete(&model.ServiceVisit{}, visit.ID)

	serviceVisitService := service.NewServiceVisitService(
		db,
		repository.NewServiceVisitRepository(db),
		repository.NewWorkOrderRepository(db),
		repository.NewWorkOrderAssignmentRepository(db),
		repository.NewTechnicianRepository(db),
		repository.NewServiceRequestRepository(db),
		repository.NewCustomerRepository(db),
		repository.NewNotificationRepository(db),
	)

	notes := "  Repair completed successfully  "
	_, err := serviceVisitService.Complete(
		visit.ID,
		technician.UserID,
		notes,
	)
	if err != nil {
		t.Fatalf("expected service visit completion to succeed, got error: %v", err)
	}

	// Pastikan Service Visit berubah menjadi COMPLETED.
	var updatedVisit model.ServiceVisit
	if err := db.First(&updatedVisit, visit.ID).Error; err != nil {
		t.Fatalf("failed to load service visit: %v", err)
	}

	if updatedVisit.Status != model.ServiceVisitCompleted {
		t.Fatalf(
			"expected visit status %v, got %v",
			model.ServiceVisitCompleted,
			updatedVisit.Status,
		)
	}

	if updatedVisit.CompletedAt == nil {
		t.Fatal("expected CompletedAt to be set")
	}

	if updatedVisit.Notes != "Repair completed successfully" {
		t.Errorf("unexpected visit notes: %q", updatedVisit.Notes)
	}

	// Work Order harus tetap IN_PROGRESS sampai laporan
	// Work Completion berhasil dibuat.
	var updatedWorkOrder model.WorkOrder
	if err := db.First(&updatedWorkOrder, workOrder.ID).Error; err != nil {
		t.Fatalf("failed to load work order: %v", err)
	}

	if updatedWorkOrder.Status != model.WorkOrderInProgress {
		t.Fatalf(
			"expected work order status %v before work completion report, got %v",
			model.WorkOrderInProgress,
			updatedWorkOrder.Status,
		)
	}

	// Notifikasi WORK_COMPLETED belum boleh dibuat pada tahap ini.
	var notifications []model.Notification
	if err := db.Where(
		"user_id = ? AND type = ? AND reference = ?",
		customer.UserID,
		model.NotificationWorkCompleted,
		fmt.Sprintf("WORK_ORDER:%d", workOrder.ID),
	).Find(&notifications).Error; err != nil {
		t.Fatalf("failed to query notifications: %v", err)
	}

	if len(notifications) != 0 {
		t.Fatalf(
			"expected no WORK_COMPLETED notification before work completion report, got %d",
			len(notifications),
		)
	}
}
