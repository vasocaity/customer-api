package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	Category  string         `json:"category" gorm:"size:50"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	// Relations
	Feedbacks []Feedback `gorm:"foreignKey:ProductID"`
}
