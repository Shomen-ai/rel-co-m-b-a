package models

import (
	"time"
)

type Specialization struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;unique"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`

	// Связи
	Doctors  []Doctor  `json:"doctors,omitempty"`
	Services []Service `json:"services,omitempty"`
}
