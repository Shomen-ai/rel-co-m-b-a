package models

import (
	"time"
)

type Doctor struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	LastName         string    `json:"last_name" gorm:"not null"`
	FirstName        string    `json:"first_name" gorm:"not null"`
	MiddleName       string    `json:"middle_name"`
	SpecializationID uint      `json:"specialization_id" gorm:"not null"`
	PhotoURL         string    `json:"photo_url"`
	ExperienceYears  int       `json:"experience_years"`
	Education        string    `json:"education"`
	IsActive         bool      `json:"is_active" gorm:"default:true"`
	DisplayOrder     int       `json:"display_order" gorm:"default:0"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Связи
	Specialization Specialization `json:"specialization,omitempty" gorm:"foreignKey:SpecializationID"`
	Services       []Service      `json:"services,omitempty" gorm:"many2many:doctor_services;"`
	Appointments   []Appointment  `json:"appointments,omitempty"`
}

// FullName возвращает полное имя врача
func (d *Doctor) FullName() string {
	if d.MiddleName != "" {
		return d.LastName + " " + d.FirstName + " " + d.MiddleName
	}
	return d.LastName + " " + d.FirstName
}
