package service

import (
	"errors"
	"strings"

	"servixa/internal/model"
	"servixa/internal/repository"
)

type TechnicianService struct {
	technicianRepository *repository.TechnicianRepository
	userRepository       *repository.UserRepository
}

func NewTechnicianService(
	technicianRepository *repository.TechnicianRepository,
	userRepository *repository.UserRepository,
) *TechnicianService {
	return &TechnicianService{
		technicianRepository: technicianRepository,
		userRepository:       userRepository,
	}
}

func (s *TechnicianService) Create(
	userID uint,
	name string,
	email string,
	phone string,
	specialization string,
) (*model.Technician, error) {

	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))
	phone = strings.TrimSpace(phone)
	specialization = strings.TrimSpace(specialization)

	if name == "" {
		return nil, errors.New("name is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	if phone == "" {
		return nil, errors.New("phone is required")
	}

	if specialization == "" {
		return nil, errors.New("specialization is required")
	}

	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Role != model.RoleTechnician {
		return nil, errors.New("only technician users can create technician profile")
	}

	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	existingTechnician, err := s.technicianRepository.FindByUserID(userID)
	if err == nil && existingTechnician != nil {
		return nil, errors.New("technician profile already exists")
	}

	technician := &model.Technician{
		UserID:         userID,
		Name:           name,
		Email:          email,
		Phone:          phone,
		Specialization: specialization,
		IsAvailable:    true,
		IsActive:       true,
	}

	if err := s.technicianRepository.Create(technician); err != nil {
		return nil, err
	}

	return technician, nil
}

func (s *TechnicianService) GetByUserID(
	userID uint,
) (*model.Technician, error) {
	technician, err := s.technicianRepository.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("technician profile not found")
	}

	if !technician.IsActive {
		return nil, errors.New("technician profile is inactive")
	}

	return technician, nil
}

func (s *TechnicianService) GetAvailable() ([]model.Technician, error) {
	return s.technicianRepository.FindAvailable()
}
