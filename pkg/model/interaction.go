package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Interaction struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerID  uuid.UUID `json:"customerId" gorm:"index; not null"`
	Channel     string    `json:"channel" gorm:"size:50"` // เช่น phone, email, chat
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
