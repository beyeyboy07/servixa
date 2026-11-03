package service

import (
	"errors"
	"strings"

	"servixa/internal/model"
	"servixa/internal/repository"

	"gorm.io/gorm"
)

type VehicleService struct {
	db                *gorm.DB
	vehicleRepository *repository.VehicleRepository
}

func NewVehicleService(
	db *gorm.DB,
	vehicleRepository *repository.VehicleRepository,
) *VehicleService {
	return &VehicleService{
		db:                db,
		vehicleRepository: vehicleRepository,
	}
}

func (s *VehicleService) Create(
	plateNumber string,
	brand string,
	modelName string,
	year int,
) (*model.Vehicle, error) {

	plateNumber = strings.ToUpper(strings.TrimSpace(plateNumber))
	brand = strings.TrimSpace(brand)
	modelName = strings.TrimSpace(modelName)

	if plateNumber == "" {
		return nil, errors.New("plate number is required")
	}

	if brand == "" {
		return nil, errors.New("brand is required")
	}

	if modelName == "" {
		return nil, errors.New("model is required")
	}

	if year <= 0 {
		return nil, errors.New("year must be greater than 0")
	}

	existingVehicle, err := s.vehicleRepository.FindByPlateNumber(plateNumber)
	if err == nil && existingVehicle != nil {
		return nil, errors.New("vehicle with plate number already exists")
	}

	vehicle := &model.Vehicle{
		PlateNumber: plateNumber,
		Brand:       brand,
		Model:       modelName,
		Year:        year,
		Status:      "AVAILABLE",
		IsActive:    true,
	}

	if err := s.vehicleRepository.Create(vehicle); err != nil {
		return nil, err
	}

	return vehicle, nil
}

func (s *VehicleService) GetByID(
	id uint,
) (*model.Vehicle, error) {
	vehicle, err := s.vehicleRepository.FindByID(id)
	if err != nil {
		return nil, errors.New("vehicle not found")
	}

	if !vehicle.IsActive {
		return nil, errors.New("vehicle is inactive")
	}

	return vehicle, nil
}

func (s *VehicleService) Update(
	id uint,
	plateNumber string,
	brand string,
	modelName string,
	year int,
	status string,
) (*model.Vehicle, error) {
	plateNumber = strings.ToUpper(strings.TrimSpace(plateNumber))
	brand = strings.TrimSpace(brand)
	modelName = strings.TrimSpace(modelName)
	status = strings.ToUpper(strings.TrimSpace(status))

	if plateNumber == "" {
		return nil, errors.New("plate number is required")
	}

	if brand == "" {
		return nil, errors.New("brand is required")
	}

	if modelName == "" {
		return nil, errors.New("model is required")
	}

	if year <= 0 {
		return nil, errors.New("year must be greater than 0")
	}

	if status == "" {
		return nil, errors.New("status is required")
	}

	vehicle, err := s.vehicleRepository.FindByID(id)
	if err != nil {
		return nil, errors.New("vehicle not found")
	}

	allowedStatuses := map[string]bool{
		"AVAILABLE":   true,
		"MAINTENANCE": true,
	}

	if !allowedStatuses[status] && status != "IN_USE" {
		return nil, errors.New("invalid vehicle status")
	}

	if vehicle.Status == "IN_USE" && status != "IN_USE" {
		return nil, errors.New("vehicle is currently assigned to a work order")
	}

	existingVehicle, err := s.vehicleRepository.FindByPlateNumber(plateNumber)
	if err == nil &&
		existingVehicle != nil &&
		existingVehicle.ID != vehicle.ID {
		return nil, errors.New("vehicle with plate number already exists")
	}

	vehicle.PlateNumber = plateNumber
	vehicle.Brand = brand
	vehicle.Model = modelName
	vehicle.Year = year
	vehicle.Status = status

	if err := s.vehicleRepository.Update(vehicle); err != nil {
		return nil, err
	}

	return vehicle, nil
}

func (s *VehicleService) GetAvailable() ([]model.Vehicle, error) {
	vehicles, err := s.vehicleRepository.FindAvailable()
	if err != nil {
		return nil, err
	}

	return vehicles, nil
}

func (s *VehicleService) GetAll() ([]model.Vehicle, error) {
	vehicles, err := s.vehicleRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return vehicles, nil
}

func (s *VehicleService) Deactivate(id uint) (*model.Vehicle, error) {
	var vehicle *model.Vehicle

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txVehicleRepository := s.vehicleRepository.WithTx(tx)

		currentVehicle, err := txVehicleRepository.FindByIDForUpdate(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("vehicle not found")
			}
			return err
		}

		if !currentVehicle.IsActive {
			return errors.New("vehicle is already inactive")
		}

		if currentVehicle.Status == "IN_USE" {
			return errors.New("vehicle is currently in use")
		}

		currentVehicle.IsActive = false

		if err := txVehicleRepository.Update(currentVehicle); err != nil {
			return err
		}

		vehicle = currentVehicle
		return nil
	})

	if err != nil {
		return nil, err
	}

	return vehicle, nil
}

func (s *VehicleService) Activate(id uint) (*model.Vehicle, error) {
	var vehicle *model.Vehicle

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txVehicleRepository := s.vehicleRepository.WithTx(tx)

		currentVehicle, err := txVehicleRepository.FindByIDForUpdate(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("vehicle not found")
			}
			return err
		}

		if currentVehicle.IsActive {
			return errors.New("vehicle is already active")
		}

		currentVehicle.IsActive = true

		if err := txVehicleRepository.Update(currentVehicle); err != nil {
			return err
		}

		vehicle = currentVehicle
		return nil
	})

	if err != nil {
		return nil, err
	}

	return vehicle, nil
}
