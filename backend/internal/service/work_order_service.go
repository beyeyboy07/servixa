package service

import (
	"errors"
	"fmt"
	"strings"

	"servixa/internal/model"
	"servixa/internal/repository"

	"gorm.io/gorm"
)

type WorkOrderService struct {
	db                            *gorm.DB
	workOrderRepository           *repository.WorkOrderRepository
	serviceRequestRepository      *repository.ServiceRequestRepository
	customerRepository            *repository.CustomerRepository
	notificationRepository        *repository.NotificationRepository
	workOrderAssignmentRepository *repository.WorkOrderAssignmentRepository
	workOrderVehicleRepository    *repository.WorkOrderVehicleRepository
	vehicleRepository             *repository.VehicleRepository
	technicianRepository          *repository.TechnicianRepository
}

func NewWorkOrderService(
	db *gorm.DB,
	workOrderRepository *repository.WorkOrderRepository,
	serviceRequestRepository *repository.ServiceRequestRepository,
	customerRepository *repository.CustomerRepository,
	notificationRepository *repository.NotificationRepository,
	workOrderAssignmentRepository *repository.WorkOrderAssignmentRepository,
	workOrderVehicleRepository *repository.WorkOrderVehicleRepository,
	vehicleRepository *repository.VehicleRepository,
	technicianRepository *repository.TechnicianRepository,
) *WorkOrderService {
	return &WorkOrderService{
		db:                            db,
		workOrderRepository:           workOrderRepository,
		serviceRequestRepository:      serviceRequestRepository,
		customerRepository:            customerRepository,
		notificationRepository:        notificationRepository,
		workOrderAssignmentRepository: workOrderAssignmentRepository,
		workOrderVehicleRepository:    workOrderVehicleRepository,
		vehicleRepository:             vehicleRepository,
		technicianRepository:          technicianRepository,
	}
}

func (s *WorkOrderService) Create(
	serviceRequestID uint,
) (*model.WorkOrder, error) {
	var workOrder *model.WorkOrder

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txServiceRequestRepository :=
			s.serviceRequestRepository.WithTx(tx)

		txWorkOrderRepository :=
			s.workOrderRepository.WithTx(tx)

		txCustomerRepository :=
			s.customerRepository.WithTx(tx)

		txNotificationRepository :=
			s.notificationRepository.WithTx(tx)

		serviceRequest, err :=
			txServiceRequestRepository.FindByIDForUpdate(
				serviceRequestID,
			)
		if err != nil {
			return errors.New("service request not found")
		}

		if serviceRequest.Status != model.ServiceRequestReviewed {
			return errors.New(
				"service request must be reviewed before creating work order",
			)
		}

		customer, err := txCustomerRepository.FindByID(
			serviceRequest.CustomerID,
		)
		if err != nil {
			return errors.New("customer not found")
		}

		existingWorkOrder, err :=
			txWorkOrderRepository.FindByServiceRequestID(
				serviceRequestID,
			)
		if err == nil && existingWorkOrder != nil {
			return errors.New("work order already exists")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		workOrder = &model.WorkOrder{
			ServiceRequestID: serviceRequest.ID,
			CustomerID:       serviceRequest.CustomerID,
			Title:            strings.TrimSpace(serviceRequest.Title),
			Description:      strings.TrimSpace(serviceRequest.Description),
			ServiceAddress:   strings.TrimSpace(serviceRequest.ServiceAddress),
			Status:           model.WorkOrderPending,
		}

		if err := txWorkOrderRepository.Create(workOrder); err != nil {
			return err
		}

		notification := &model.Notification{
			UserID:    customer.UserID,
			Type:      model.NotificationWorkOrderCreated,
			Title:     "Work Order Created",
			Message:   fmt.Sprintf("Your work order has been created: %s", workOrder.Title),
			Reference: fmt.Sprintf("WORK_ORDER:%d", workOrder.ID),
			IsRead:    false,
		}

		if err := txNotificationRepository.Create(notification); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return workOrder, nil
}

func (s *WorkOrderService) Cancel(
	workOrderID uint,
) (*model.WorkOrder, error) {
	var workOrder *model.WorkOrder

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txWorkOrderRepository :=
			s.workOrderRepository.WithTx(tx)

		txCustomerRepository :=
			s.customerRepository.WithTx(tx)

		txAssignmentRepository :=
			s.workOrderAssignmentRepository.WithTx(tx)

		txWorkOrderVehicleRepository :=
			s.workOrderVehicleRepository.WithTx(tx)

		txTechnicianRepository :=
			s.technicianRepository.WithTx(tx)

		txNotificationRepository :=
			s.notificationRepository.WithTx(tx)

		txVehicleRepository := s.vehicleRepository.WithTx(tx)

		var err error

		workOrder, err =
			txWorkOrderRepository.FindByIDForUpdate(workOrderID)
		if err != nil {
			return errors.New("work order not found")
		}

		switch workOrder.Status {
		case model.WorkOrderCompleted,
			model.WorkOrderCustomerConfirmed,
			model.WorkOrderClosed,
			model.WorkOrderCancelled:
			return errors.New("work order cannot be cancelled")
		}

		customer, err :=
			txCustomerRepository.FindByID(workOrder.CustomerID)
		if err != nil {
			return errors.New("customer not found")
		}

		// Update the status inside the same transaction.
		workOrder.Status = model.WorkOrderCancelled

		if err := txWorkOrderRepository.Update(workOrder); err != nil {
			return err
		}

		// Release assigned vehicle, if any.
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

		reference := fmt.Sprintf("WORK_ORDER:%d", workOrder.ID)

		// Notify the customer.
		customerNotification := &model.Notification{
			UserID:    customer.UserID,
			Type:      model.NotificationWorkOrderCancelled,
			Title:     "Work Order Cancelled",
			Message:   fmt.Sprintf("Your work order has been cancelled: %s", workOrder.Title),
			Reference: reference,
			IsRead:    false,
		}

		if err := txNotificationRepository.Create(customerNotification); err != nil {
			return err
		}

		// Notify the assigned technician, if an assignment exists.
		assignment, err :=
			txAssignmentRepository.FindByWorkOrderID(workOrder.ID)

		if err == nil && assignment != nil {
			technician, err :=
				txTechnicianRepository.FindByID(assignment.TechnicianID)
			if err != nil {
				return errors.New("assigned technician not found")
			}

			technicianNotification := &model.Notification{
				UserID:    technician.UserID,
				Type:      model.NotificationWorkOrderCancelled,
				Title:     "Work Order Cancelled",
				Message:   fmt.Sprintf("Your assigned work order has been cancelled: %s", workOrder.Title),
				Reference: reference,
				IsRead:    false,
			}

			if err := txNotificationRepository.Create(technicianNotification); err != nil {
				return err
			}
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return workOrder, nil
}

func (s *WorkOrderService) GetByID(
	id uint,
) (
	*model.WorkOrder,
	*model.WorkOrderAssignment,
	*model.WorkOrderVehicle,
	error,
) {
	workOrder, err := s.workOrderRepository.FindByID(id)
	if err != nil {
		return nil, nil, nil, errors.New("work order not found")
	}

	assignment, err :=
		s.workOrderAssignmentRepository.FindByWorkOrderID(workOrder.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			assignment = nil
		} else {
			return nil, nil, nil, err
		}
	}

	vehicleAssignment, err :=
		s.workOrderVehicleRepository.FindByWorkOrderID(workOrder.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			vehicleAssignment = nil
		} else {
			return nil, nil, nil, err
		}
	}

	return workOrder, assignment, vehicleAssignment, nil
}

func (s *WorkOrderService) GetAll(
	userID uint,
	userRole model.UserRole,
) ([]model.WorkOrder, error) {
	switch userRole {
	case model.RoleAdmin, model.RoleDispatcher:
		return s.workOrderRepository.FindAll()

	case model.RoleTechnician:
		technician, err := s.technicianRepository.FindByUserID(userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return []model.WorkOrder{}, nil
			}
			return nil, err
		}

		assignments, err :=
			s.workOrderAssignmentRepository.FindByTechnicianID(technician.ID)
		if err != nil {
			return nil, err
		}

		workOrders := make([]model.WorkOrder, 0, len(assignments))

		for _, assignment := range assignments {
			workOrder, err :=
				s.workOrderRepository.FindByID(assignment.WorkOrderID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					continue
				}
				return nil, err
			}

			workOrders = append(workOrders, *workOrder)
		}

		return workOrders, nil

	case model.RoleCustomer:
		customer, err := s.customerRepository.FindByUserID(userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return []model.WorkOrder{}, nil
			}
			return nil, err
		}

		return s.workOrderRepository.FindByCustomerID(customer.ID)

	default:
		return nil, errors.New("invalid user role")
	}
}
