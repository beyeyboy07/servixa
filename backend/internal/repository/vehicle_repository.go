package repository

import (
	"servixa/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VehicleRepository struct {
	*BaseRepository
}

func NewVehicleRepository(db *gorm.DB) *VehicleRepository {
	return &VehicleRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *VehicleRepository) Create(vehicle *model.Vehicle) error {
	return r.db.Create(vehicle).Error
}

func (r *VehicleRepository) FindByID(id uint) (*model.Vehicle, error) {
	var vehicle model.Vehicle

	err := r.db.
		Where("id = ?", id).
		First(&vehicle).Error

	if err != nil {
		return nil, err
	}

	return &vehicle, nil
}

func (r *VehicleRepository) FindByPlateNumber(
	plateNumber string,
) (*model.Vehicle, error) {
	var vehicle model.Vehicle

	err := r.db.
		Where("plate_number = ?", plateNumber).
		First(&vehicle).Error

	if err != nil {
		return nil, err
	}

	return &vehicle, nil
}

func (r *VehicleRepository) Update(vehicle *model.Vehicle) error {
	return r.db.Save(vehicle).Error
}

func (r *VehicleRepository) FindByIDForUpdate(
	id uint,
) (*model.Vehicle, error) {
	var vehicle model.Vehicle

	err := r.db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&vehicle).
		Error

	if err != nil {
		return nil, err
	}

	return &vehicle, nil
}

func (r *VehicleRepository) WithTx(
	tx *gorm.DB,
) *VehicleRepository {
	return NewVehicleRepository(tx)
}

func (r *VehicleRepository) FindAvailable() ([]model.Vehicle, error) {
	var vehicles []model.Vehicle

	err := r.db.
		Where("is_active = ? AND status = ?", true, "AVAILABLE").
		Order("plate_number ASC").
		Find(&vehicles).Error

	if err != nil {
		return nil, err
	}

	return vehicles, nil
}

func (r *VehicleRepository) FindAll() ([]model.Vehicle, error) {
	var vehicles []model.Vehicle

	err := r.db.
		Order("plate_number ASC").
		Find(&vehicles).Error

	if err != nil {
		return nil, err
	}

	return vehicles, nil
}
