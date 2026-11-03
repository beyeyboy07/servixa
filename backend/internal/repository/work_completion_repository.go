package repository

import (
	"servixa/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WorkCompletionRepository struct {
	*BaseRepository
}

func NewWorkCompletionRepository(db *gorm.DB) *WorkCompletionRepository {
	return &WorkCompletionRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *WorkCompletionRepository) Create(
	completion *model.WorkCompletion,
) error {
	return r.db.Create(completion).Error
}

func (r *WorkCompletionRepository) FindByID(
	id uint,
) (*model.WorkCompletion, error) {
	var completion model.WorkCompletion

	err := r.db.
		Where("id = ?", id).
		First(&completion).
		Error

	if err != nil {
		return nil, err
	}

	return &completion, nil
}

func (r *WorkCompletionRepository) FindByWorkOrderID(
	workOrderID uint,
) (*model.WorkCompletion, error) {
	var completion model.WorkCompletion

	err := r.db.
		Where("work_order_id = ?", workOrderID).
		First(&completion).
		Error

	if err != nil {
		return nil, err
	}

	return &completion, nil
}

func (r *WorkCompletionRepository) FindByWorkOrderIDForUpdate(
	workOrderID uint,
) (*model.WorkCompletion, error) {
	var completion model.WorkCompletion

	err := r.db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("work_order_id = ?", workOrderID).
		First(&completion).
		Error

	if err != nil {
		return nil, err
	}

	return &completion, nil
}

func (r *WorkCompletionRepository) Update(
	completion *model.WorkCompletion,
) error {
	return r.db.Save(completion).Error
}

func (r *WorkCompletionRepository) WithTx(
	tx *gorm.DB,
) *WorkCompletionRepository {
	return NewWorkCompletionRepository(tx)
}
