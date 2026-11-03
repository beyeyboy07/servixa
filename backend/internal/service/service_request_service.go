package service

import (
	"errors"
	"strings"
	"time"

	"servixa/internal/model"
	"servixa/internal/repository"
)

type ServiceRequestService struct {
	serviceRequestRepository *repository.ServiceRequestRepository
	customerRepository       *repository.CustomerRepository
	userRepository           *repository.UserRepository
}

func NewServiceRequestService(
	serviceRequestRepository *repository.ServiceRequestRepository,
	customerRepository *repository.CustomerRepository,
	userRepository *repository.UserRepository,
) *ServiceRequestService {
	return &ServiceRequestService{
		serviceRequestRepository: serviceRequestRepository,
		customerRepository:       customerRepository,
		userRepository:           userRepository,
	}
}

func (s *ServiceRequestService) Create(
	userID uint,
	title string,
	description string,
	serviceAddress string,
	preferredDate *time.Time,
) (*model.ServiceRequest, error) {

	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)
	serviceAddress = strings.TrimSpace(serviceAddress)

	if title == "" {
		return nil, errors.New("title is required")
	}

	if description == "" {
		return nil, errors.New("description is required")
	}

	if serviceAddress == "" {
		return nil, errors.New("service address is required")
	}

	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Role != model.RoleCustomer {
		return nil, errors.New("only customer users can create service request")
	}

	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	customer, err := s.customerRepository.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("customer profile not found")
	}

	if !customer.IsActive {
		return nil, errors.New("customer profile is inactive")
	}

	serviceRequest := &model.ServiceRequest{
		CustomerID:     customer.ID,
		Title:          title,
		Description:    description,
		ServiceAddress: serviceAddress,
		PreferredDate:  preferredDate,
		Status:         model.ServiceRequestPending,
	}

	if err := s.serviceRequestRepository.Create(serviceRequest); err != nil {
		return nil, err
	}

	return serviceRequest, nil
}

func (s *ServiceRequestService) GetByID(
	userID uint,
	id uint,
) (*model.ServiceRequest, error) {

	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Role != model.RoleCustomer {
		return nil, errors.New("only customer users can access service request")
	}

	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	customer, err := s.customerRepository.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("customer profile not found")
	}

	if !customer.IsActive {
		return nil, errors.New("customer profile is inactive")
	}

	serviceRequest, err := s.serviceRequestRepository.FindByID(id)
	if err != nil {
		return nil, errors.New("service request not found")
	}

	if serviceRequest.CustomerID != customer.ID {
		return nil, errors.New("service request not found")
	}

	return serviceRequest, nil
}

func (s *ServiceRequestService) GetByCustomer(
	userID uint,
) ([]model.ServiceRequest, error) {

	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Role != model.RoleCustomer {
		return nil, errors.New("only customer users can access service requests")
	}

	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	customer, err := s.customerRepository.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("customer profile not found")
	}

	if !customer.IsActive {
		return nil, errors.New("customer profile is inactive")
	}

	serviceRequests, err := s.serviceRequestRepository.FindByCustomerID(customer.ID)
	if err != nil {
		return nil, err
	}

	return serviceRequests, nil
}

func (s *ServiceRequestService) Review(
	id uint,
) (*model.ServiceRequest, error) {

	serviceRequest, err := s.serviceRequestRepository.FindByID(id)
	if err != nil {
		return nil, errors.New("service request not found")
	}

	if serviceRequest.Status != model.ServiceRequestPending {
		return nil, errors.New("service request cannot be reviewed")
	}

	serviceRequest.Status = model.ServiceRequestReviewed

	if err := s.serviceRequestRepository.Update(serviceRequest); err != nil {
		return nil, err
	}

	return serviceRequest, nil
}

func (s *ServiceRequestService) Reject(
	id uint,
) (*model.ServiceRequest, error) {

	serviceRequest, err := s.serviceRequestRepository.FindByID(id)
	if err != nil {
		return nil, errors.New("service request not found")
	}

	if serviceRequest.Status != model.ServiceRequestPending {
		return nil, errors.New("service request cannot be rejected")
	}

	serviceRequest.Status = model.ServiceRequestRejected

	if err := s.serviceRequestRepository.Update(serviceRequest); err != nil {
		return nil, err
	}

	return serviceRequest, nil
}

func (s *ServiceRequestService) GetAll(
	userID uint,
) ([]model.ServiceRequest, error) {

	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	if user.Role != model.RoleAdmin && user.Role != model.RoleDispatcher {
		return nil, errors.New("only admin and dispatcher users can access service requests")
	}

	return s.serviceRequestRepository.FindAll()
}
