package repository

import (
	"customer-api/pkg/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(fd *model.Product) error
	GetByID(id uuid.UUID) (*model.Product, error)
	Update(cus *model.Product) error
	Delete(id uuid.UUID) error
	List(query string, limit, offset int) ([]model.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

// Create implements ProductRepository.
func (f *productRepository) Create(fd *model.Product) error {
	return f.db.Create(fd).Error
}

// Delete implements ProductRepository.
func (f *productRepository) Delete(id uuid.UUID) error {
	return f.db.Delete(id).Error
}

// GetByID implements ProductRepository.
func (f *productRepository) GetByID(id uuid.UUID) (*model.Product, error) {
	var Product model.Product

	if err := f.db.First(&Product, id).Error; err != nil {
		return nil, err
	}

	return &Product, nil
}

// List implements ProductRepository.
func (f *productRepository) List(query string, limit int, offset int) ([]model.Product, error) {
	var list []model.Product

	result := f.db.Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}

// Update implements ProductRepository.
func (f *productRepository) Update(fd *model.Product) error {
	return f.db.Save(fd).Error
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}
