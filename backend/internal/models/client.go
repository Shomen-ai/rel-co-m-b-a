package models

import (
	"time"
)

type Client struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Phone      string `json:"phone" gorm:"not null;index"`
	LastName   string `json:"last_name" gorm:"not null"`
	FirstName  string `json:"first_name" gorm:"not null"`
	MiddleName string `json:"middle_name"`
	BirthDate  string `json:"birth_date"` // YYYY-MM-DD

	ContactPhone    bool `json:"contact_phone" gorm:"default:false"`
	ContactTelegram bool `json:"contact_telegram" gorm:"default:false"`
	ContactWhatsApp bool `json:"contact_whatsapp" gorm:"default:false"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Связи
	Appointments []Appointment `json:"appointments,omitempty" gorm:"foreignKey:ClientID"`
}

// FullName возвращает полное имя клиента
func (c *Client) FullName() string {
	if c.MiddleName != "" {
		return c.LastName + " " + c.FirstName + " " + c.MiddleName
	}
	return c.LastName + " " + c.FirstName
}
