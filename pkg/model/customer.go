package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string         `json:"name" gorm:"size:255;not null"`
	Email     string         `json:"email" gorm:"size:255; uniqueIndex;not null"`
	Phone     string         `json:"phone" gorm:"size:50"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
