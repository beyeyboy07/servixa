package repository

import (
	"servixa/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ServiceRequestRepository struct {
	*BaseRepository
}

func NewServiceRequestRepository(db *gorm.DB) *ServiceRequestRepository {
	return &ServiceRequestRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *ServiceRequestRepository) Create(
	serviceRequest *model.ServiceRequest,
) error {
	return r.db.Create(serviceRequest).Error
}

func (r *ServiceRequestRepository) FindByID(
	id uint,
) (*model.ServiceRequest, error) {

	var serviceRequest model.ServiceRequest

	err := r.db.
		Where("id = ?", id).
		First(&serviceRequest).Error

	if err != nil {
		return nil, err
	}

	return &serviceRequest, nil
}

func (r *ServiceRequestRepository) FindByCustomerID(
	customerID uint,
) ([]model.ServiceRequest, error) {

	var serviceRequests []model.ServiceRequest

	err := r.db.
		Where("customer_id = ?", customerID).
		Order("created_at DESC").
		Find(&serviceRequests).Error

	if err != nil {
		return nil, err
	}

	return serviceRequests, nil
}

func (r *ServiceRequestRepository) Update(
	serviceRequest *model.ServiceRequest,
) error {
	return r.db.Save(serviceRequest).Error
}

func (r *ServiceRequestRepository) WithTx(tx *gorm.DB) *ServiceRequestRepository {
	return NewServiceRequestRepository(tx)
}

func (r *ServiceRequestRepository) FindByIDForUpdate(
	id uint,
) (*model.ServiceRequest, error) {
	var serviceRequest model.ServiceRequest

	err := r.db.
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where("id = ?", id).
		First(&serviceRequest).
		Error

	if err != nil {
		return nil, err
	}

	return &serviceRequest, nil
}

func (r *ServiceRequestRepository) FindAll() ([]model.ServiceRequest, error) {
	var serviceRequests []model.ServiceRequest

	err := r.db.
		Order("created_at DESC").
		Find(&serviceRequests).Error

	if err != nil {
		return nil, err
	}

	return serviceRequests, nil
}
