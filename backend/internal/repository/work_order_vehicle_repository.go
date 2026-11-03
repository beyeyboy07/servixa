package repository

import (
	"servixa/internal/model"

	"gorm.io/gorm"
)

type WorkOrderVehicleRepository struct {
	*BaseRepository
}

func NewWorkOrderVehicleRepository(
	db *gorm.DB,
) *WorkOrderVehicleRepository {
	return &WorkOrderVehicleRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *WorkOrderVehicleRepository) Create(
	workOrderVehicle *model.WorkOrderVehicle,
) error {
	return r.db.Create(workOrderVehicle).Error
}

func (r *WorkOrderVehicleRepository) FindByWorkOrderID(
	workOrderID uint,
) (*model.WorkOrderVehicle, error) {
	var workOrderVehicle model.WorkOrderVehicle

	err := r.db.
		Where("work_order_id = ?", workOrderID).
		First(&workOrderVehicle).
		Error

	if err != nil {
		return nil, err
	}

	return &workOrderVehicle, nil
}

func (r *WorkOrderVehicleRepository) FindByVehicleID(
	vehicleID uint,
) ([]model.WorkOrderVehicle, error) {
	var workOrderVehicles []model.WorkOrderVehicle

	err := r.db.
		Where("vehicle_id = ?", vehicleID).
		Order("created_at DESC").
		Find(&workOrderVehicles).
		Error

	if err != nil {
		return nil, err
	}

	return workOrderVehicles, nil
}

func (r *WorkOrderVehicleRepository) WithTx(
	tx *gorm.DB,
) *WorkOrderVehicleRepository {
	return NewWorkOrderVehicleRepository(tx)
}
