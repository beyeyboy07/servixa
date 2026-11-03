package service

import (
	"errors"
	"strings"

	"servixa/internal/model"
	"servixa/internal/repository"
)

type PartService struct {
	partRepository *repository.PartRepository
}

func NewPartService(
	partRepository *repository.PartRepository,
) *PartService {
	return &PartService{
		partRepository: partRepository,
	}
}

func (s *PartService) Create(
	code string,
	name string,
	description string,
	unit string,
	stock int,
	minStock int,
	unitPrice float64,
) (*model.Part, error) {

	code = strings.ToUpper(strings.TrimSpace(code))
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	unit = strings.ToLower(strings.TrimSpace(unit))

	if code == "" {
		return nil, errors.New("code is required")
	}

	if name == "" {
		return nil, errors.New("name is required")
	}

	if unit == "" {
		return nil, errors.New("unit is required")
	}

	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}

	if minStock < 0 {
		return nil, errors.New("minimum stock cannot be negative")
	}

	if unitPrice < 0 {
		return nil, errors.New("unit price cannot be negative")
	}

	existingPart, err := s.partRepository.FindByCode(code)
	if err == nil && existingPart != nil {
		return nil, errors.New("part with code already exists")
	}

	part := &model.Part{
		Code:        code,
		Name:        name,
		Description: description,
		Unit:        unit,
		Stock:       stock,
		MinStock:    minStock,
		UnitPrice:   unitPrice,
		IsActive:    true,
	}

	if err := s.partRepository.Create(part); err != nil {
		return nil, err
	}

	return part, nil
}

func (s *PartService) GetByID(
	id uint,
) (*model.Part, error) {
	part, err := s.partRepository.FindByID(id)
	if err != nil {
		return nil, errors.New("part not found")
	}

	if !part.IsActive {
		return nil, errors.New("part is inactive")
	}

	return part, nil
}

func (s *PartService) Update(
	id uint,
	code string,
	name string,
	description string,
	unit string,
	stock int,
	minStock int,
	unitPrice float64,
) (*model.Part, error) {

	code = strings.ToUpper(strings.TrimSpace(code))
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	unit = strings.ToLower(strings.TrimSpace(unit))

	if code == "" {
		return nil, errors.New("code is required")
	}

	if name == "" {
		return nil, errors.New("name is required")
	}

	if unit == "" {
		return nil, errors.New("unit is required")
	}

	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}

	if minStock < 0 {
		return nil, errors.New("minimum stock cannot be negative")
	}

	if unitPrice < 0 {
		return nil, errors.New("unit price cannot be negative")
	}

	part, err := s.partRepository.FindByID(id)
	if err != nil {
		return nil, errors.New("part not found")
	}

	if !part.IsActive {
		return nil, errors.New("part is inactive")
	}

	existingPart, err := s.partRepository.FindByCode(code)
	if err == nil && existingPart != nil && existingPart.ID != part.ID {
		return nil, errors.New("part with code already exists")
	}

	part.Code = code
	part.Name = name
	part.Description = description
	part.Unit = unit
	part.Stock = stock
	part.MinStock = minStock
	part.UnitPrice = unitPrice

	if err := s.partRepository.Update(part); err != nil {
		return nil, err
	}

	return part, nil
}
