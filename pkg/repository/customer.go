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

type repository struct {
	db *gorm.DB
}

// Create implements CustomerRepository.
func (r *repository) Create(cus *model.Customer) error {
	return r.db.Create(cus).Error
}

// Delete implements CustomerRepository.
func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Customer{}, id).Error
}

// GetByID implements CustomerRepository.
func (r *repository) GetByID(id uuid.UUID) (*model.Customer, error) {
	var c model.Customer
	if err := r.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// List implements CustomerRepository.
func (r *repository) List(query string, limit int, offset int) ([]model.Customer, error) {
	var list []model.Customer
	err := r.db.
		Order("name asc").
		Limit(limit).
		Offset(offset).
		Find(&list, "name LIKE '%?%' OR email LIKE '%?%'", query, query)

	if err != nil {
		return nil, err.Error
	}
	return list, nil
}

// Update implements CustomerRepository.
func (r *repository) Update(cus *model.Customer) error {
	return r.db.Save(cus).Error
}

func NewRepository(db *gorm.DB) CustomerRepository {
	return &repository{db: db}
}
