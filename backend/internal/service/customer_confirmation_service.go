package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"servixa/internal/model"
	"servixa/internal/repository"

	"gorm.io/gorm"
)

type CustomerConfirmationService struct {
	db                             *gorm.DB
	customerConfirmationRepository *repository.CustomerConfirmationRepository
	workOrderRepository            *repository.WorkOrderRepository
	workCompletionRepository       *repository.WorkCompletionRepository
	customerRepository             *repository.CustomerRepository
	workOrderAssignmentRepository  *repository.WorkOrderAssignmentRepository
	technicianRepository           *repository.TechnicianRepository
	notificationRepository         *repository.NotificationRepository
}

func NewCustomerConfirmationService(
	db *gorm.DB,
	customerConfirmationRepository *repository.CustomerConfirmationRepository,
	workOrderRepository *repository.WorkOrderRepository,
	workCompletionRepository *repository.WorkCompletionRepository,
	customerRepository *repository.CustomerRepository,
	workOrderAssignmentRepository *repository.WorkOrderAssignmentRepository,
	technicianRepository *repository.TechnicianRepository,
	notificationRepository *repository.NotificationRepository,
) *CustomerConfirmationService {
	return &CustomerConfirmationService{
		db:                             db,
		customerConfirmationRepository: customerConfirmationRepository,
		workOrderRepository:            workOrderRepository,
		workCompletionRepository:       workCompletionRepository,
		customerRepository:             customerRepository,
		workOrderAssignmentRepository:  workOrderAssignmentRepository,
		technicianRepository:           technicianRepository,
		notificationRepository:         notificationRepository,
	}
}

func (s *CustomerConfirmationService) Confirm(
	workOrderID uint,
	userID uint,
	status model.CustomerConfirmationStatus,
	notes string,
) (*model.CustomerConfirmation, error) {
	var confirmation *model.CustomerConfirmation

	if status != model.CustomerConfirmationConfirmed &&
		status != model.CustomerConfirmationRejected {
		return nil, errors.New("invalid customer confirmation status")
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txCustomerRepository := s.customerRepository.WithTx(tx)
		txWorkOrderRepository := s.workOrderRepository.WithTx(tx)
		txWorkCompletionRepository := s.workCompletionRepository.WithTx(tx)
		txConfirmationRepository := s.customerConfirmationRepository.WithTx(tx)
		txAssignmentRepository := s.workOrderAssignmentRepository.WithTx(tx)
		txTechnicianRepository := s.technicianRepository.WithTx(tx)
		txNotificationRepository := s.notificationRepository.WithTx(tx)

		customer, err := txCustomerRepository.FindByUserID(userID)
		if err != nil {
			return errors.New("customer not found")
		}

		if !customer.IsActive {
			return errors.New("customer is inactive")
		}

		workOrder, err := txWorkOrderRepository.FindByIDForUpdate(workOrderID)
		if err != nil {
			return errors.New("work order not found")
		}

		if workOrder.CustomerID != customer.ID {
			return errors.New("customer is not authorized for this work order")
		}

		if workOrder.Status != model.WorkOrderCompleted {
			return errors.New("work order must be completed before confirmation")
		}

		_, err = txWorkCompletionRepository.FindByWorkOrderID(workOrderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("work completion not found")
			}
			return err
		}

		existing, err := txConfirmationRepository.FindByWorkOrderIDForUpdate(workOrderID)
		if err == nil && existing != nil {
			return errors.New("customer confirmation already exists")
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		confirmation = &model.CustomerConfirmation{
			WorkOrderID: workOrder.ID,
			CustomerID:  customer.ID,
			Status:      status,
			Notes:       strings.TrimSpace(notes),
		}

		if status == model.CustomerConfirmationConfirmed {
			now := time.Now()
			confirmation.ConfirmedAt = &now
		}

		if err := txConfirmationRepository.Create(confirmation); err != nil {
			return err
		}

		if status == model.CustomerConfirmationConfirmed {
			workOrder.Status = model.WorkOrderClosed

			if err := txWorkOrderRepository.Update(workOrder); err != nil {
				return err
			}

			assignment, err := txAssignmentRepository.FindByWorkOrderID(workOrder.ID)
			if err != nil {
				return errors.New("technician assignment not found")
			}

			technician, err := txTechnicianRepository.FindByID(assignment.TechnicianID)
			if err != nil {
				return errors.New("technician not found")
			}

			notification := &model.Notification{
				UserID: technician.UserID,
				Type:   model.NotificationCustomerConfirmed,
				Title:  "Customer Confirmed Work",
				Message: fmt.Sprintf(
					"Customer confirmed the completed work order: %s",
					workOrder.Title,
				),
				Reference: fmt.Sprintf("WORK_ORDER:%d", workOrder.ID),
				IsRead:    false,
			}

			if err := txNotificationRepository.Create(notification); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return confirmation, nil
}
