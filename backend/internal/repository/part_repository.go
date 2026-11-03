package repository

import (
	"servixa/internal/model"

	"gorm.io/gorm"
)

type PartRepository struct {
	*BaseRepository
}

func NewPartRepository(db *gorm.DB) *PartRepository {
	return &PartRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *PartRepository) Create(part *model.Part) error {
	return r.db.Create(part).Error
}

func (r *PartRepository) FindByID(id uint) (*model.Part, error) {
	var part model.Part

	err := r.db.
		Where("id = ?", id).
		First(&part).Error

	if err != nil {
		return nil, err
	}

	return &part, nil
}

func (r *PartRepository) FindByCode(code string) (*model.Part, error) {
	var part model.Part

	err := r.db.
		Where("code = ?", code).
		First(&part).Error

	if err != nil {
		return nil, err
	}

	return &part, nil
}

func (r *PartRepository) Update(part *model.Part) error {
	return r.db.Save(part).Error
}
