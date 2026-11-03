package repository

import (
	"servixa/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TechnicianRepository struct {
	*BaseRepository
}

func NewTechnicianRepository(db *gorm.DB) *TechnicianRepository {
	return &TechnicianRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *TechnicianRepository) Create(technician *model.Technician) error {
	return r.db.Create(technician).Error
}

func (r *TechnicianRepository) FindByID(id uint) (*model.Technician, error) {
	var technician model.Technician

	err := r.db.
		Where("id = ?", id).
		First(&technician).Error

	if err != nil {
		return nil, err
	}

	return &technician, nil
}

func (r *TechnicianRepository) FindByUserID(userID uint) (*model.Technician, error) {
	var technician model.Technician

	err := r.db.
		Where("user_id = ?", userID).
		First(&technician).Error

	if err != nil {
		return nil, err
	}

	return &technician, nil
}

func (r *TechnicianRepository) Update(technician *model.Technician) error {
	return r.db.Save(technician).Error
}

func (r *TechnicianRepository) FindByIDForUpdate(
	id uint,
) (*model.Technician, error) {
	var technician model.Technician

	err := r.db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&technician).
		Error

	if err != nil {
		return nil, err
	}

	return &technician, nil
}

func (r *TechnicianRepository) WithTx(
	tx *gorm.DB,
) *TechnicianRepository {
	return NewTechnicianRepository(tx)
}

func (r *TechnicianRepository) FindAvailable() ([]model.Technician, error) {
	var technicians []model.Technician

	err := r.db.
		Where("is_active = ? AND is_available = ?", true, true).
		Order("name ASC").
		Find(&technicians).Error

	if err != nil {
		return nil, err
	}

	return technicians, nil
}
