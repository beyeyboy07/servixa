package repository

import (
	"servixa/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WorkOrderRepository struct {
	*BaseRepository
}

func NewWorkOrderRepository(db *gorm.DB) *WorkOrderRepository {
	return &WorkOrderRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *WorkOrderRepository) Create(
	workOrder *model.WorkOrder,
) error {
	return r.db.Create(workOrder).Error
}

func (r *WorkOrderRepository) FindByID(
	id uint,
) (*model.WorkOrder, error) {
	var workOrder model.WorkOrder

	err := r.db.
		Where("id = ?", id).
		First(&workOrder).
		Error

	if err != nil {
		return nil, err
	}

	return &workOrder, nil
}

func (r *WorkOrderRepository) FindByServiceRequestID(
	serviceRequestID uint,
) (*model.WorkOrder, error) {
	var workOrder model.WorkOrder

	err := r.db.
		Where("service_request_id = ?", serviceRequestID).
		First(&workOrder).
		Error

	if err != nil {
		return nil, err
	}

	return &workOrder, nil
}

func (r *WorkOrderRepository) FindByCustomerID(
	customerID uint,
) ([]model.WorkOrder, error) {
	var workOrders []model.WorkOrder

	err := r.db.
		Where("customer_id = ?", customerID).
		Order("created_at DESC").
		Find(&workOrders).
		Error

	if err != nil {
		return nil, err
	}

	return workOrders, nil
}

func (r *WorkOrderRepository) Update(
	workOrder *model.WorkOrder,
) error {
	return r.db.Save(workOrder).Error
}

func (r *WorkOrderRepository) FindAll() ([]model.WorkOrder, error) {
	var workOrders []model.WorkOrder

	err := r.db.
		Order("created_at DESC").
		Find(&workOrders).
		Error

	if err != nil {
		return nil, err
	}

	return workOrders, nil
}

func (r *WorkOrderRepository) WithTx(tx *gorm.DB) *WorkOrderRepository {
	return NewWorkOrderRepository(tx)
}

func (r *WorkOrderRepository) FindByIDForUpdate(
	id uint,
) (*model.WorkOrder, error) {
	var workOrder model.WorkOrder

	err := r.db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&workOrder).
		Error

	if err != nil {
		return nil, err
	}

	return &workOrder, nil
}
