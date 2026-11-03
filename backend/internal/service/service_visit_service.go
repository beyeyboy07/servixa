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

type ServiceVisitService struct {
	db                       *gorm.DB
	serviceVisitRepository   *repository.ServiceVisitRepository
	workOrderRepository      *repository.WorkOrderRepository
	workOrderAssignmentRepo  *repository.WorkOrderAssignmentRepository
	technicianRepository     *repository.TechnicianRepository
	serviceRequestRepository *repository.ServiceRequestRepository
	customerRepository       *repository.CustomerRepository
	notificationRepository   *repository.NotificationRepository
}

func NewServiceVisitService(
	db *gorm.DB,
	serviceVisitRepository *repository.ServiceVisitRepository,
	workOrderRepository *repository.WorkOrderRepository,
	workOrderAssignmentRepo *repository.WorkOrderAssignmentRepository,
	technicianRepository *repository.TechnicianRepository,
	serviceRequestRepository *repository.ServiceRequestRepository,
	customerRepository *repository.CustomerRepository,
	notificationRepository *repository.NotificationRepository,
) *ServiceVisitService {
	return &ServiceVisitService{
		db:                       db,
		serviceVisitRepository:   serviceVisitRepository,
		workOrderRepository:      workOrderRepository,
		workOrderAssignmentRepo:  workOrderAssignmentRepo,
		technicianRepository:     technicianRepository,
		serviceRequestRepository: serviceRequestRepository,
		customerRepository:       customerRepository,
		notificationRepository:   notificationRepository,
	}
}

func (s *ServiceVisitService) Create(
	workOrderID uint,
	userID uint,
) (*model.ServiceVisit, error) {
	var visit *model.ServiceVisit

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txWorkOrderRepository := s.workOrderRepository.WithTx(tx)
		txAssignmentRepository := s.workOrderAssignmentRepo.WithTx(tx)
		txVisitRepository := s.serviceVisitRepository.WithTx(tx)
		txTechnicianRepository := s.technicianRepository.WithTx(tx)

		technician, err := txTechnicianRepository.FindByUserID(userID)
		if err != nil {
			return errors.New("technician not found")
		}

		if !technician.IsActive {
			return errors.New("technician is inactive")
		}

		workOrder, err := txWorkOrderRepository.FindByIDForUpdate(workOrderID)
		if err != nil {
			return errors.New("work order not found")
		}

		if workOrder.Status != model.WorkOrderScheduled {
			return errors.New("work order must be scheduled")
		}

		assignment, err := txAssignmentRepository.FindByWorkOrderID(workOrderID)
		if err != nil {
			return errors.New("technician assignment not found")
		}

		if assignment.TechnicianID != technician.ID {
			return errors.New("technician is not assigned to this work order")
		}

		existingVisit, err := txVisitRepository.FindByWorkOrderID(workOrderID)

		if err == nil && existingVisit != nil {
			return errors.New("service visit already exists")
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		visit = &model.ServiceVisit{
			WorkOrderID:  workOrderID,
			TechnicianID: technician.ID,
			Status:       model.ServiceVisitScheduled,
		}

		if err := txVisitRepository.Create(visit); err != nil {
			return err
		}

		workOrder.Status = model.WorkOrderTechnicianAccepted

		if err := txWorkOrderRepository.Update(workOrder); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return visit, nil
}

func (s *ServiceVisitService) Start(
	visitID uint,
	userID uint,
) (*model.ServiceVisit, error) {
	var visit *model.ServiceVisit

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txVisitRepository := s.serviceVisitRepository.WithTx(tx)
		txWorkOrderRepository := s.workOrderRepository.WithTx(tx)
		txTechnicianRepository := s.technicianRepository.WithTx(tx)
		txServiceRequestRepository := s.serviceRequestRepository.WithTx(tx)
		txCustomerRepository := s.customerRepository.WithTx(tx)
		txNotificationRepository := s.notificationRepository.WithTx(tx)

		technician, err := txTechnicianRepository.FindByUserID(userID)
		if err != nil {
			return errors.New("technician not found")
		}

		if !technician.IsActive {
			return errors.New("technician is inactive")
		}

		visit, err = txVisitRepository.FindByIDForUpdate(visitID)
		if err != nil {
			return errors.New("service visit not found")
		}

		if visit.TechnicianID != technician.ID {
			return errors.New("technician is not assigned to this visit")
		}

		if visit.Status != model.ServiceVisitScheduled {
			return errors.New("service visit cannot be started")
		}

		workOrder, err := txWorkOrderRepository.FindByIDForUpdate(visit.WorkOrderID)
		if err != nil {
			return errors.New("work order not found")
		}

		if workOrder.Status != model.WorkOrderTechnicianAccepted {
			return errors.New("work order cannot be started")
		}

		serviceRequest, err := txServiceRequestRepository.FindByID(workOrder.ServiceRequestID)
		if err != nil {
			return errors.New("service request not found")
		}

		customer, err := txCustomerRepository.FindByID(serviceRequest.CustomerID)
		if err != nil {
			return errors.New("customer not found")
		}

		now := time.Now()

		visit.Status = model.ServiceVisitInProgress
		visit.StartedAt = &now

		if err := txVisitRepository.Update(visit); err != nil {
			return err
		}

		workOrder.Status = model.WorkOrderInProgress

		if err := txWorkOrderRepository.Update(workOrder); err != nil {
			return err
		}

		notification := &model.Notification{
			UserID:  customer.UserID,
			Type:    model.NotificationWorkStarted,
			Title:   "Work Started",
			Message: fmt.Sprintf("Work has started on your work order: %s", workOrder.Title),
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

	return visit, nil
}

func (s *ServiceVisitService) Complete(
	visitID uint,
	userID uint,
	notes string,
) (*model.ServiceVisit, error) {
	var visit *model.ServiceVisit

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txVisitRepository := s.serviceVisitRepository.WithTx(tx)
		txWorkOrderRepository := s.workOrderRepository.WithTx(tx)
		txTechnicianRepository := s.technicianRepository.WithTx(tx)

		technician, err := txTechnicianRepository.FindByUserID(userID)
		if err != nil {
			return errors.New("technician not found")
		}

		if !technician.IsActive {
			return errors.New("technician is inactive")
		}

		visit, err = txVisitRepository.FindByIDForUpdate(visitID)
		if err != nil {
			return errors.New("service visit not found")
		}

		if visit.TechnicianID != technician.ID {
			return errors.New("technician is not assigned to this visit")
		}

		if visit.Status != model.ServiceVisitInProgress {
			return errors.New("service visit cannot be completed")
		}

		workOrder, err := txWorkOrderRepository.FindByIDForUpdate(
			visit.WorkOrderID,
		)
		if err != nil {
			return errors.New("work order not found")
		}

		if workOrder.Status != model.WorkOrderInProgress {
			return errors.New("work order cannot be completed")
		}

		now := time.Now()

		visit.Status = model.ServiceVisitCompleted
		visit.CompletedAt = &now
		visit.Notes = strings.TrimSpace(notes)

		if err := txVisitRepository.Update(visit); err != nil {
			return err
		}

		// Work Order tetap IN_PROGRESS.
		// Status COMPLETED baru ditetapkan setelah Work Completion dibuat.
		return nil
	})

	if err != nil {
		return nil, err
	}

	return visit, nil
}
