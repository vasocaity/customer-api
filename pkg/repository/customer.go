package repository

import (
	"customer-api/pkg/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(cus *model.Customer) error
	GetByID(id uuid.UUID) (*model.Customer, error)
	Update(cus *model.Customer) error
	Delete(id uuid.UUID) error
	List(query string, limit, offset int) ([]model.Customer, error)
}

type customerRepository struct {
	db *gorm.DB
}

// Create implements CustomerRepository.
func (r *customerRepository) Create(cus *model.Customer) error {
	return r.db.Create(cus).Error
}

// Delete implements CustomerRepository.
func (r *customerRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Customer{}, id).Error
}

// GetByID implements CustomerRepository.
func (r *customerRepository) GetByID(id uuid.UUID) (*model.Customer, error) {
	var c model.Customer
	if err := r.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// List implements CustomerRepository.
func (r *customerRepository) List(query string, limit int, offset int) ([]model.Customer, error) {
	var list []model.Customer
	result := r.db.
		Order("name asc").
		Limit(limit).
		Offset(offset).
		Preload("Feedbacks").
		Preload("Interactions").
		Find(&list)

	if result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}

// Update implements CustomerRepository.
func (r *customerRepository) Update(cus *model.Customer) error {
	return r.db.Save(cus).Error
}

func NewRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}
