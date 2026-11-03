package repository

import (
	"servixa/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ServiceVisitRepository struct {
	*BaseRepository
}

func NewServiceVisitRepository(
	db *gorm.DB,
) *ServiceVisitRepository {
	return &ServiceVisitRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *ServiceVisitRepository) Create(
	visit *model.ServiceVisit,
) error {
	return r.db.Create(visit).Error
}

func (r *ServiceVisitRepository) FindByID(
	id uint,
) (*model.ServiceVisit, error) {
	var visit model.ServiceVisit

	err := r.db.
		Where("id = ?", id).
		First(&visit).
		Error

	if err != nil {
		return nil, err
	}

	return &visit, nil
}

func (r *ServiceVisitRepository) FindByWorkOrderID(
	workOrderID uint,
) (*model.ServiceVisit, error) {
	var visit model.ServiceVisit

	err := r.db.
		Where("work_order_id = ?", workOrderID).
		First(&visit).
		Error

	if err != nil {
		return nil, err
	}

	return &visit, nil
}

func (r *ServiceVisitRepository) Update(
	visit *model.ServiceVisit,
) error {
	return r.db.Save(visit).Error
}

func (r *ServiceVisitRepository) FindByIDForUpdate(
	id uint,
) (*model.ServiceVisit, error) {
	var visit model.ServiceVisit

	err := r.db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&visit).
		Error

	if err != nil {
		return nil, err
	}

	return &visit, nil
}

func (r *ServiceVisitRepository) WithTx(
	tx *gorm.DB,
) *ServiceVisitRepository {
	return NewServiceVisitRepository(tx)
}
