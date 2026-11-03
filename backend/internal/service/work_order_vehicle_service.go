package service

import (
	"errors"
	"time"

	"servixa/internal/model"
	"servixa/internal/repository"

	"gorm.io/gorm"
)

type WorkOrderVehicleService struct {
	db                         *gorm.DB
	workOrderVehicleRepository *repository.WorkOrderVehicleRepository
	workOrderRepository        *repository.WorkOrderRepository
	vehicleRepository          *repository.VehicleRepository
}

func NewWorkOrderVehicleService(
	db *gorm.DB,
	workOrderVehicleRepository *repository.WorkOrderVehicleRepository,
	workOrderRepository *repository.WorkOrderRepository,
	vehicleRepository *repository.VehicleRepository,
) *WorkOrderVehicleService {
	return &WorkOrderVehicleService{
		db:                         db,
		workOrderVehicleRepository: workOrderVehicleRepository,
		workOrderRepository:        workOrderRepository,
		vehicleRepository:          vehicleRepository,
	}
}

func (s *WorkOrderVehicleService) Assign(
	workOrderID uint,
	vehicleID uint,
) (*model.WorkOrderVehicle, error) {

	var assignment *model.WorkOrderVehicle

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txWorkOrderRepository := s.workOrderRepository.WithTx(tx)
		txVehicleRepository := s.vehicleRepository.WithTx(tx)
		txAssignmentRepository := s.workOrderVehicleRepository.WithTx(tx)

		workOrder, err := txWorkOrderRepository.FindByIDForUpdate(workOrderID)
		if err != nil {
			return errors.New("work order not found")
		}

		if workOrder.Status == model.WorkOrderCompleted ||
			workOrder.Status == model.WorkOrderCustomerConfirmed ||
			workOrder.Status == model.WorkOrderClosed ||
			workOrder.Status == model.WorkOrderCancelled {
			return errors.New("work order cannot be assigned")
		}

		vehicle, err := txVehicleRepository.FindByIDForUpdate(vehicleID)
		if err != nil {
			return errors.New("vehicle not found")
		}

		if !vehicle.IsActive {
			return errors.New("vehicle is inactive")
		}

		if vehicle.Status != "AVAILABLE" {
			return errors.New("vehicle is not available")
		}

		existingAssignment, err :=
			txAssignmentRepository.FindByWorkOrderID(workOrderID)

		if err == nil && existingAssignment != nil {
			return errors.New("work order already has a vehicle")
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		assignment = &model.WorkOrderVehicle{
			WorkOrderID: workOrderID,
			VehicleID:   vehicleID,
			AssignedAt:  time.Now(),
		}

		if err := txAssignmentRepository.Create(assignment); err != nil {
			return err
		}

		vehicle.Status = "IN_USE"

		if err := txVehicleRepository.Update(vehicle); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return assignment, nil
}
