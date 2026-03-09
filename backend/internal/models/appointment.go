// internal/models/appointment.go
package models

import (
	"time"
)

// TimeSlot - слот времени
type TimeSlot struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	DoctorID      uint      `json:"doctor_id" gorm:"not null;index"`
	SlotDate      string    `json:"slot_date" gorm:"not null"`
	StartTime     string    `json:"start_time" gorm:"not null"`
	EndTime       string    `json:"end_time" gorm:"not null"`
	IsAvailable   bool      `json:"is_available" gorm:"default:true"`
	AppointmentID *uint     `json:"appointment_id" gorm:"unique"` // unique constraint
	CreatedAt     time.Time `json:"created_at"`

	// Связи
	Doctor      Doctor       `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`
	Appointment *Appointment `json:"appointment,omitempty" gorm:"foreignKey:AppointmentID"`
}

// Appointment - запись на приём
type Appointment struct {
	ID         uint `json:"id" gorm:"primaryKey"`
	ClientID   uint `json:"client_id" gorm:"not null;index"`
	DoctorID   uint `json:"doctor_id" gorm:"not null;index"`
	TimeSlotID uint `json:"time_slot_id" gorm:"not null;unique"`

	Status    string    `json:"status" gorm:"default:'pending'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Связи
	Client   Client   `json:"client,omitempty" gorm:"foreignKey:ClientID"`
	Doctor   Doctor   `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`
	TimeSlot TimeSlot `json:"time_slot,omitempty" gorm:"foreignKey:TimeSlotID"`
}

// DoctorWorkDay - рабочий день врача (для статистики)
type DoctorWorkDay struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	DoctorID  uint      `json:"doctor_id" gorm:"not null"`
	Doctor    Doctor    `json:"doctor,omitempty"`
	WorkDate  string    `json:"work_date" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// Request/Response структуры
type AvailableSlotsResponse struct {
	Date  string     `json:"date"`
	Slots []TimeSlot `json:"slots"`
}

type CreateAppointmentRequest struct {
	DoctorID        uint   `json:"doctor_id" binding:"required"`
	TimeSlotID      uint   `json:"time_slot_id" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	FirstName       string `json:"first_name" binding:"required"`
	MiddleName      string `json:"middle_name"`
	BirthDate       string `json:"birth_date" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
	ContactPhone    bool   `json:"contact_phone"`
	ContactTelegram bool   `json:"contact_telegram"`
	ContactWhatsApp bool   `json:"contact_whatsapp"`
}
