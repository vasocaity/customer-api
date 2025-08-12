package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Feedback struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerID uuid.UUID `json:"customerId" gorm:"index; not null"`
	ProductID  uuid.UUID `json:"productId" gorm:"index;not null"`
	Rating     int       `gorm:"not null"` // 1-5
	Comment    string    `gorm:"type:text"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
