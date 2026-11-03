package tests

import (
	"fmt"
	"testing"

	"servixa/internal/config"
	"servixa/internal/database"
	"servixa/internal/model"
	"servixa/internal/repository"
	"servixa/internal/service"
)

func TestWorkOrderAssignmentRejectsInvalidTechnician(t *testing.T) {
	testCases := []struct {
		name        string
		isActive    bool
		isAvailable bool
		expectedErr string
	}{
		{
			name:        "inactive technician",
			isActive:    false,
			isAvailable: true,
			expectedErr: "technician is inactive",
		},
		{
			name:        "unavailable technician",
			isActive:    true,
			isAvailable: false,
			expectedErr: "technician is not available",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.Load()
			db := database.Connect(cfg)

			customer := &model.Customer{
				UserID:   1,
				Name:     "Assignment Validation Customer",
				Email:    "assignment-validation@servixa.local",
				Phone:    "0800000001",
				Address:  "Assignment Validation Address",
				IsActive: true,
			}

			if err := db.Create(customer).Error; err != nil {
				t.Fatalf("failed to create customer: %v", err)
			}
			defer db.Delete(&model.Customer{}, customer.ID)

			workOrder := &model.WorkOrder{
				ServiceRequestID: 999999,
				CustomerID:       customer.ID,
				Title:            "Assignment Validation Test",
				Description:      "Test invalid technician assignment",
				ServiceAddress:   customer.Address,
				Status:           model.WorkOrderPending,
			}

			if err := db.Create(workOrder).Error; err != nil {
				t.Fatalf("failed to create work order: %v", err)
			}
			defer db.Delete(&model.WorkOrder{}, workOrder.ID)

			technician := &model.Technician{
				UserID:         1,
				Name:           "Assignment Validation Technician",
				Email:          "assignment-validation-tech@servixa.local",
				Phone:          "0800000002",
				Specialization: "General Service",
				IsAvailable:    true,
				IsActive:       true,
			}

			result := db.Exec(`
	INSERT INTO technicians (
		user_id,
		name,
		email,
		phone,
		specialization,
		is_available,
		is_active,
		created_at,
		updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
`,
				technician.UserID,
				technician.Name,
				technician.Email,
				technician.Phone,
				technician.Specialization,
				true,
				true,
			)

			if result.Error != nil {
				t.Fatalf("failed to create technician: %v", result.Error)
			}

			if err := db.Raw(
				"SELECT id FROM technicians WHERE email = ?",
				technician.Email,
			).Scan(&technician.ID).Error; err != nil {
				t.Fatalf("failed to fetch technician ID: %v", err)
			}

			defer db.Exec("DELETE FROM technicians WHERE id = ?", technician.ID)

			if err := db.Exec(
				"UPDATE technicians SET is_active = ?, is_available = ? WHERE id = ?",
				tc.isActive,
				tc.isAvailable,
				technician.ID,
			).Error; err != nil {
				t.Fatalf("failed to update technician status: %v", err)
			}

			assignmentRepository :=
				repository.NewWorkOrderAssignmentRepository(db)
			workOrderRepository :=
				repository.NewWorkOrderRepository(db)
			technicianRepository :=
				repository.NewTechnicianRepository(db)
			notificationRepository :=
				repository.NewNotificationRepository(db)

			assignmentService :=
				service.NewWorkOrderAssignmentService(
					db,
					assignmentRepository,
					workOrderRepository,
					technicianRepository,
					notificationRepository,
				)

			_, err := assignmentService.Assign(
				workOrder.ID,
				technician.ID,
			)

			if err == nil {
				t.Fatal("expected assignment to be rejected")
			}

			if err.Error() != tc.expectedErr {
				t.Fatalf(
					"expected error %q, got %q",
					tc.expectedErr,
					err.Error(),
				)
			}

			var assignmentCount int64
			if err := db.Model(&model.WorkOrderAssignment{}).
				Where("work_order_id = ?", workOrder.ID).
				Count(&assignmentCount).Error; err != nil {
				t.Fatalf("failed to count assignments: %v", err)
			}

			if assignmentCount != 0 {
				t.Fatalf(
					"expected no assignments, got %d",
					assignmentCount,
				)
			}

			var updatedWorkOrder model.WorkOrder
			if err := db.First(
				&updatedWorkOrder,
				workOrder.ID,
			).Error; err != nil {
				t.Fatalf("failed to fetch work order: %v", err)
			}

			if updatedWorkOrder.Status != model.WorkOrderPending {
				t.Fatalf(
					"expected status %s, got %s",
					model.WorkOrderPending,
					updatedWorkOrder.Status,
				)
			}

			var notificationCount int64
			if err := db.Model(&model.Notification{}).
				Where(
					"user_id = ? AND reference = ?",
					technician.UserID,
					fmt.Sprintf("WORK_ORDER:%d", workOrder.ID),
				).
				Count(&notificationCount).Error; err != nil {
				t.Fatalf("failed to count notifications: %v", err)
			}

			if notificationCount != 0 {
				t.Fatalf(
					"expected no assignment notifications, got %d",
					notificationCount,
				)
			}
		})
	}
}
