package tests

import (
	"testing"

	"servixa/internal/config"
	"servixa/internal/database"
	"servixa/internal/model"
	"servixa/internal/repository"
	"servixa/internal/service"
)

func TestWorkCompletionTechnicianAuthorization(t *testing.T) {
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
		Name:     "Work Completion Authorization Customer",
		Email:    "work-completion-auth@servixa.local",
		Phone:    "0800000010",
		Address:  "Test Address",
		IsActive: true,
	}

	if err := db.Create(customer).Error; err != nil {
		t.Fatalf("failed to create test customer: %v", err)
	}
	defer db.Delete(&model.Customer{}, customer.ID)

	serviceRequest := &model.ServiceRequest{
		CustomerID:     customer.ID,
		Title:          "Authorization Service Request",
		Description:    "Authorization test",
		ServiceAddress: "Test Address",
		Status:         model.ServiceRequestReviewed,
	}

	if err := db.Create(serviceRequest).Error; err != nil {
		t.Fatalf("failed to create service request: %v", err)
	}
	defer db.Delete(&model.ServiceRequest{}, serviceRequest.ID)

	workOrder := &model.WorkOrder{
		ServiceRequestID: serviceRequest.ID,
		CustomerID:       customer.ID,
		Title:            "Work Completion Authorization Test",
		Description:      "Authorization test",
		ServiceAddress:   "Test Address",
		Status:           model.WorkOrderInProgress,
	}

	if err := db.Create(workOrder).Error; err != nil {
		t.Fatalf("failed to create test work order: %v", err)
	}
	defer db.Delete(&model.WorkOrder{}, workOrder.ID)

	assignedTechnician := &model.Technician{
		UserID:         999991,
		Name:           "Assigned Technician",
		Email:          "assigned-completion@servixa.local",
		Phone:          "0800000011",
		Specialization: "General Service",
		IsAvailable:    true,
		IsActive:       true,
	}

	if err := db.Create(assignedTechnician).Error; err != nil {
		t.Fatalf("failed to create assigned technician: %v", err)
	}
	defer db.Delete(&model.Technician{}, assignedTechnician.ID)

	otherTechnician := &model.Technician{
		UserID:         999992,
		Name:           "Other Technician",
		Email:          "other-completion@servixa.local",
		Phone:          "0800000012",
		Specialization: "General Service",
		IsAvailable:    true,
		IsActive:       true,
	}

	if err := db.Create(otherTechnician).Error; err != nil {
		t.Fatalf("failed to create other technician: %v", err)
	}
	defer db.Delete(&model.Technician{}, otherTechnician.ID)

	assignment := &model.WorkOrderAssignment{
		WorkOrderID:  workOrder.ID,
		TechnicianID: assignedTechnician.ID,
	}

	if err := db.Create(assignment).Error; err != nil {
		t.Fatalf("failed to create technician assignment: %v", err)
	}
	defer db.Delete(&model.WorkOrderAssignment{}, assignment.ID)

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

	_, err := workCompletionService.Create(
		workOrder.ID,
		otherTechnician.UserID,
		"Unauthorized completion",
		"Trying to complete another technician's work order",
		"",
	)

	if err == nil {
		t.Fatal("expected unauthorized technician to be rejected")
	}

	expectedError := "technician is not assigned to this work order"

	if err.Error() != expectedError {
		t.Fatalf(
			"expected error %q, got %q",
			expectedError,
			err.Error(),
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

	if completionCount != 0 {
		t.Fatalf(
			"expected no work completion to be created, got %d",
			completionCount,
		)
	}
}
