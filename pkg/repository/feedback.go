package repository

import (
	"customer-api/pkg/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FeedbackRepository interface {
	Create(fd *model.Feedback) error
	GetByID(id uuid.UUID) (*model.Feedback, error)
	Update(cus *model.Feedback) error
	Delete(id uuid.UUID) error
	List(query string, limit, offset int) ([]model.Feedback, error)
}

type feedbackRepository struct {
	db *gorm.DB
}

// Create implements FeedbackRepository.
func (f *feedbackRepository) Create(fd *model.Feedback) error {
	return f.db.Create(fd).Error
}

// Delete implements FeedbackRepository.
func (f *feedbackRepository) Delete(id uuid.UUID) error {
	return f.db.Delete(id).Error
}

// GetByID implements FeedbackRepository.
func (f *feedbackRepository) GetByID(id uuid.UUID) (*model.Feedback, error) {
	var feedback model.Feedback

	if err := f.db.First(&feedback, id).Error; err != nil {
		return nil, err
	}

	return &feedback, nil
}

// List implements FeedbackRepository.
func (f *feedbackRepository) List(query string, limit int, offset int) ([]model.Feedback, error) {
	var list []model.Feedback

	result := f.db.Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}

// Update implements FeedbackRepository.
func (f *feedbackRepository) Update(fd *model.Feedback) error {
	return f.db.Save(fd).Error
}

func NewFeedbackRepository(db *gorm.DB) FeedbackRepository {
	return &feedbackRepository{
		db: db,
	}
}
