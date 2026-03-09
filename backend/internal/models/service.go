package models

import (
	"time"
)

type Service struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Name             string    `json:"name" gorm:"not null"`
	Code             string    `json:"code"` // Код услуги (например, А11.05.001)
	Description      string    `json:"description"`
	Price            float64   `json:"price" gorm:"not null"`
	DurationMinutes  int       `json:"duration_minutes"`
	SpecializationID uint      `json:"specialization_id" gorm:"not null"`
	IsActive         bool      `json:"is_active" gorm:"default:true"`
	DisplayOrder     int       `json:"display_order" gorm:"default:0"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Связи
	Specialization Specialization `json:"specialization,omitempty" gorm:"foreignKey:SpecializationID"`
	Doctors        []Doctor       `json:"doctors,omitempty" gorm:"many2many:doctor_services;"`
}

// Junction table для связи многие-ко-многим между Doctor и Service
type DoctorService struct {
	DoctorID  uint      `gorm:"primaryKey"`
	ServiceID uint      `gorm:"primaryKey"`
	Price     float64   `json:"price"` // Можно указать цену конкретного врага на услугу
	CreatedAt time.Time `json:"created_at"`
}
