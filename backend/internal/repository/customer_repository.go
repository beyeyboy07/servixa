package repository

import (
	"servixa/internal/model"

	"gorm.io/gorm"
)

type CustomerRepository struct {
	*BaseRepository
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *CustomerRepository) Create(customer *model.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepository) FindByID(id uint) (*model.Customer, error) {
	var customer model.Customer

	err := r.db.
		Where("id = ?", id).
		First(&customer).Error

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *CustomerRepository) FindByUserID(userID uint) (*model.Customer, error) {
	var customer model.Customer

	err := r.db.
		Where("user_id = ?", userID).
		First(&customer).Error

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *CustomerRepository) Update(customer *model.Customer) error {
	return r.db.Save(customer).Error
}

func (r *CustomerRepository) WithTx(tx *gorm.DB) *CustomerRepository {
	return NewCustomerRepository(tx)
}

func (r *CustomerRepository) FindAll() ([]model.Customer, error) {
	var customers []model.Customer

	err := r.db.
		Order("created_at DESC").
		Find(&customers).Error

	if err != nil {
		return nil, err
	}

	return customers, nil
}
