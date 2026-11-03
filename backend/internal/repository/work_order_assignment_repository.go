package repository

import (
	"servixa/internal/model"

	"gorm.io/gorm"
)

type WorkOrderAssignmentRepository struct {
	*BaseRepository
}

func NewWorkOrderAssignmentRepository(
	db *gorm.DB,
) *WorkOrderAssignmentRepository {
	return &WorkOrderAssignmentRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *WorkOrderAssignmentRepository) Create(
	assignment *model.WorkOrderAssignment,
) error {
	return r.db.Create(assignment).Error
}

func (r *WorkOrderAssignmentRepository) FindByWorkOrderID(
	workOrderID uint,
) (*model.WorkOrderAssignment, error) {
	var assignment model.WorkOrderAssignment

	err := r.db.
		Where("work_order_id = ?", workOrderID).
		First(&assignment).
		Error

	if err != nil {
		return nil, err
	}

	return &assignment, nil
}

func (r *WorkOrderAssignmentRepository) FindByTechnicianID(
	technicianID uint,
) ([]model.WorkOrderAssignment, error) {
	var assignments []model.WorkOrderAssignment

	err := r.db.
		Where("technician_id = ?", technicianID).
		Order("created_at DESC").
		Find(&assignments).
		Error

	if err != nil {
		return nil, err
	}

	return assignments, nil
}

func (r *WorkOrderAssignmentRepository) WithTx(
	tx *gorm.DB,
) *WorkOrderAssignmentRepository {
	return NewWorkOrderAssignmentRepository(tx)
}

func (r *WorkOrderAssignmentRepository) Update(
	assignment *model.WorkOrderAssignment,
) error {
	return r.db.Save(assignment).Error
}
