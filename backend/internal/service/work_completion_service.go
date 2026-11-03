package service

import (
	"errors"
	"fmt"
	"strings"

	"servixa/internal/model"
	"servixa/internal/repository"

	"gorm.io/gorm"
)

type WorkCompletionService struct {
	db                         *gorm.DB
	workCompletionRepository   *repository.WorkCompletionRepository
	workOrderRepository        *repository.WorkOrderRepository
	technicianRepository       *repository.TechnicianRepository
	workOrderAssignmentRepo    *repository.WorkOrderAssignmentRepository
	serviceVisitRepository     *repository.ServiceVisitRepository
	serviceRequestRepository   *repository.ServiceRequestRepository
	customerRepository         *repository.CustomerRepository
	notificationRepository     *repository.NotificationRepository
	workOrderVehicleRepository *repository.WorkOrderVehicleRepository
	vehicleRepository          *repository.VehicleRepository
}

func NewWorkCompletionService(
	db *gorm.DB,
	workCompletionRepository *repository.WorkCompletionRepository,
	workOrderRepository *repository.WorkOrderRepository,
	workOrderVehicleRepository *repository.WorkOrderVehicleRepository,
	vehicleRepository *repository.VehicleRepository,
	workOrderAssignmentRepo *repository.WorkOrderAssignmentRepository,
	technicianRepository *repository.TechnicianRepository,
	serviceVisitRepository *repository.ServiceVisitRepository,
	serviceRequestRepository *repository.ServiceRequestRepository,
	customerRepository *repository.CustomerRepository,
	notificationRepository *repository.NotificationRepository,
) *WorkCompletionService {
	return &WorkCompletionService{
		db:                         db,
		workCompletionRepository:   workCompletionRepository,
		workOrderRepository:        workOrderRepository,
		workOrderVehicleRepository: workOrderVehicleRepository,
		vehicleRepository:          vehicleRepository,
		workOrderAssignmentRepo:    workOrderAssignmentRepo,
		technicianRepository:       technicianRepository,
		serviceVisitRepository:     serviceVisitRepository,
		serviceRequestRepository:   serviceRequestRepository,
		customerRepository:         customerRepository,
		notificationRepository:     notificationRepository,
	}
}

func (s *WorkCompletionService) Create(
	workOrderID uint,
	userID uint,
	summary string,
	workPerformed string,
	recommendations string,
) (*model.WorkCompletion, error) {
	var completion *model.WorkCompletion

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txWorkOrderRepository := s.workOrderRepository.WithTx(tx)

		txWorkOrderVehicleRepository := s.workOrderVehicleRepository.WithTx(tx)

		txTechnicianRepository := s.technicianRepository.WithTx(tx)
		txWorkCompletionRepository := s.workCompletionRepository.WithTx(tx)
		txAssignmentRepository := s.workOrderAssignmentRepo.WithTx(tx)
		txVisitRepository := s.serviceVisitRepository.WithTx(tx)
		txServiceRequestRepository := s.serviceRequestRepository.WithTx(tx)
		txCustomerRepository := s.customerRepository.WithTx(tx)
		txNotificationRepository := s.notificationRepository.WithTx(tx)
		txVehicleRepository := s.vehicleRepository.WithTx(tx)

		technician, err := txTechnicianRepository.FindByUserID(userID)
		if err != nil {
			return errors.New("technician not found")
		}

		if !technician.IsActive {
			return errors.New("technician is inactive")
		}

		// Lock Work Order to serialize concurrent completion requests.
		workOrder, err := txWorkOrderRepository.FindByIDForUpdate(workOrderID)
		if err != nil {
			return errors.New("work order not found")
		}

		// Service Visit must be completed before its report can be submitted.
		if workOrder.Status != model.WorkOrderInProgress {
			return errors.New("work order must be in progress")
		}

		assignment, err := txAssignmentRepository.FindByWorkOrderID(workOrderID)
		if err != nil {
			return errors.New("technician assignment not found")
		}

		if assignment.TechnicianID != technician.ID {
			return errors.New("technician is not assigned to this work order")
		}

		visit, err := txVisitRepository.FindByWorkOrderID(workOrderID)
		if err != nil {
			return errors.New("service visit not found")
		}

		if visit.TechnicianID != technician.ID {
			return errors.New("technician is not assigned to this visit")
		}

		if visit.Status != model.ServiceVisitCompleted {
			return errors.New("service visit must be completed first")
		}

		// Prevent duplicate reports.
		existingCompletion, err :=
			txWorkCompletionRepository.FindByWorkOrderIDForUpdate(workOrderID)

		if err == nil && existingCompletion != nil {
			return errors.New("work completion already exists")
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		summary = strings.TrimSpace(summary)
		workPerformed = strings.TrimSpace(workPerformed)
		recommendations = strings.TrimSpace(recommendations)

		if summary == "" {
			return errors.New("summary is required")
		}

		if workPerformed == "" {
			return errors.New("work performed is required")
		}

		completion = &model.WorkCompletion{
			WorkOrderID:     workOrderID,
			TechnicianID:    technician.ID,
			Summary:         summary,
			WorkPerformed:   workPerformed,
			Recommendations: recommendations,
		}

		if err := txWorkCompletionRepository.Create(completion); err != nil {
			return err
		}

		// Only mark the Work Order completed after the report is saved.
		workOrder.Status = model.WorkOrderCompleted

		if err := txWorkOrderRepository.Update(workOrder); err != nil {
			return err
		}

		// Release assigned vehicle after Work Order is completed.
		// Release assigned vehicle after Work Order is completed.
		vehicleAssignment, err :=
			txWorkOrderVehicleRepository.FindByWorkOrderID(workOrder.ID)

		if err == nil && vehicleAssignment != nil {
			vehicle, err := txVehicleRepository.FindByIDForUpdate(
				vehicleAssignment.VehicleID,
			)
			if err != nil {
				return err
			}

			if vehicle.Status != "IN_USE" {
				return errors.New("assigned vehicle is not marked as in use")
			}

			vehicle.Status = "AVAILABLE"

			if err := txVehicleRepository.Update(vehicle); err != nil {
				return err
			}

		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// Find the customer who owns the original Service Request.
		serviceRequest, err := txServiceRequestRepository.FindByID(
			workOrder.ServiceRequestID,
		)
		if err != nil {
			return errors.New("service request not found")
		}

		customer, err := txCustomerRepository.FindByID(
			serviceRequest.CustomerID,
		)
		if err != nil {
			return errors.New("customer not found")
		}

		notification := &model.Notification{
			UserID:  customer.UserID,
			Type:    model.NotificationWorkCompleted,
			Title:   "Work Completed",
			Message: fmt.Sprintf("Your work order has been completed: %s", workOrder.Title),
			Reference: fmt.Sprintf(
				"WORK_ORDER:%d",
				workOrder.ID,
			),
			IsRead: false,
		}

		if err := txNotificationRepository.Create(notification); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return completion, nil
}
