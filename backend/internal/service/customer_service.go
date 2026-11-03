package service

import (
	"errors"
	"strings"

	"servixa/internal/model"
	"servixa/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CustomerService struct {
	db                 *gorm.DB
	customerRepository *repository.CustomerRepository
	userRepository     *repository.UserRepository
}

func NewCustomerService(
	db *gorm.DB,
	customerRepository *repository.CustomerRepository,
	userRepository *repository.UserRepository,
) *CustomerService {
	return &CustomerService{
		db:                 db,
		customerRepository: customerRepository,
		userRepository:     userRepository,
	}
}

func (s *CustomerService) Create(
	userID uint,
	name string,
	email string,
	phone string,
	address string,
) (*model.Customer, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))
	phone = strings.TrimSpace(phone)
	address = strings.TrimSpace(address)

	if name == "" {
		return nil, errors.New("name is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	if phone == "" {
		return nil, errors.New("phone is required")
	}

	if address == "" {
		return nil, errors.New("address is required")
	}

	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Role != model.RoleCustomer {
		return nil, errors.New("only customer users can create customer profile")
	}

	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	existingCustomer, err := s.customerRepository.FindByUserID(userID)
	if err == nil && existingCustomer != nil {
		return nil, errors.New("customer profile already exists")
	}

	customer := &model.Customer{
		UserID:   userID,
		Name:     name,
		Email:    email,
		Phone:    phone,
		Address:  address,
		IsActive: true,
	}

	if err := s.customerRepository.Create(customer); err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *CustomerService) GetByUserID(userID uint) (*model.Customer, error) {
	customer, err := s.customerRepository.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("customer profile not found")
	}

	if !customer.IsActive {
		return nil, errors.New("customer profile is inactive")
	}

	return customer, nil
}

func (s *CustomerService) GetAll() ([]model.Customer, error) {
	return s.customerRepository.FindAll()
}

func (s *CustomerService) GetByID(id uint) (*model.Customer, error) {
	customer, err := s.customerRepository.FindByID(id)
	if err != nil {
		return nil, errors.New("customer not found")
	}

	return customer, nil
}

func (s *CustomerService) CreateByAdmin(
	name string,
	email string,
	password string,
	phone string,
	address string,
) (*model.Customer, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))
	password = strings.TrimSpace(password)
	phone = strings.TrimSpace(phone)
	address = strings.TrimSpace(address)

	if name == "" {
		return nil, errors.New("name is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	if password == "" {
		return nil, errors.New("password is required")
	}

	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	if phone == "" {
		return nil, errors.New("phone is required")
	}

	if address == "" {
		return nil, errors.New("address is required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	var customer *model.Customer

	err = s.db.Transaction(func(tx *gorm.DB) error {
		txUserRepository := s.userRepository.WithTx(tx)
		txCustomerRepository := s.customerRepository.WithTx(tx)

		existingUser, err := txUserRepository.FindByEmail(email)

		if err == nil && existingUser != nil {
			return errors.New("email already registered")
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		user := &model.User{
			Name:     name,
			Email:    email,
			Password: string(hashedPassword),
			Role:     model.RoleCustomer,
			IsActive: true,
		}

		if err := txUserRepository.Create(user); err != nil {
			return err
		}

		customer = &model.Customer{
			UserID:   user.ID,
			Name:     name,
			Email:    email,
			Phone:    phone,
			Address:  address,
			IsActive: true,
		}

		if err := txCustomerRepository.Create(customer); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return customer, nil
}
