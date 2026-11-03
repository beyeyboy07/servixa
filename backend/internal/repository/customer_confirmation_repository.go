package repository

import (
	"servixa/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CustomerConfirmationRepository struct {
	*BaseRepository
}

func NewCustomerConfirmationRepository(
	db *gorm.DB,
) *CustomerConfirmationRepository {
	return &CustomerConfirmationRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *CustomerConfirmationRepository) Create(
	confirmation *model.CustomerConfirmation,
) error {
	return r.db.Create(confirmation).Error
}

func (r *CustomerConfirmationRepository) FindByID(
	id uint,
) (*model.CustomerConfirmation, error) {
	var confirmation model.CustomerConfirmation

	err := r.db.
		Where("id = ?", id).
		First(&confirmation).
		Error

	if err != nil {
		return nil, err
	}

	return &confirmation, nil
}

func (r *CustomerConfirmationRepository) FindByWorkOrderID(
	workOrderID uint,
) (*model.CustomerConfirmation, error) {
	var confirmation model.CustomerConfirmation

	err := r.db.
		Where("work_order_id = ?", workOrderID).
		First(&confirmation).
		Error

	if err != nil {
		return nil, err
	}

	return &confirmation, nil
}

func (r *CustomerConfirmationRepository) FindByWorkOrderIDForUpdate(
	workOrderID uint,
) (*model.CustomerConfirmation, error) {
	var confirmation model.CustomerConfirmation

	err := r.db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("work_order_id = ?", workOrderID).
		First(&confirmation).
		Error

	if err != nil {
		return nil, err
	}

	return &confirmation, nil
}

func (r *CustomerConfirmationRepository) WithTx(
	tx *gorm.DB,
) *CustomerConfirmationRepository {
	return NewCustomerConfirmationRepository(tx)
}
