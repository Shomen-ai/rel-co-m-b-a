// internal/handlers/appointment_handlers.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"healthcare/internal/database"
	"healthcare/internal/models"

	"github.com/gorilla/mux"
)

// Request/Response структуры
type CreateAppointmentRequest struct {
	// Данные клиента (нового или существующего)
	ClientID   *uint  `json:"client_id"` // если клиент уже существует
	Phone      string `json:"phone" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	FirstName  string `json:"first_name" binding:"required"`
	MiddleName string `json:"middle_name"`
	BirthDate  string `json:"birth_date" binding:"required"`

	// Способы связи
	ContactPhone    bool `json:"contact_phone"`
	ContactTelegram bool `json:"contact_telegram"`
	ContactWhatsApp bool `json:"contact_whatsapp"`

	// Данные записи
	DoctorID   uint `json:"doctor_id" binding:"required"`
	TimeSlotID uint `json:"time_slot_id" binding:"required"`
}

type AppointmentResponse struct {
	ID        uint            `json:"id"`
	Client    models.Client   `json:"client"`
	Doctor    models.Doctor   `json:"doctor"`
	TimeSlot  models.TimeSlot `json:"time_slot"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
}

// GetAvailableSlots - получить свободные слоты на дату
func GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	var slots []models.TimeSlot

	result := database.DB.
		Where("slot_date = ? AND is_available = ? AND appointment_id IS NULL", date, true).
		Preload("Doctor").
		Preload("Doctor.Specialization").
		Order("start_time asc").
		Find(&slots)

	if result.Error != nil {
		http.Error(w, "Ошибка получения слотов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(slots)
}

// GetDoctorSlots - получить слоты конкретного врача на дату
func GetDoctorSlots(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	doctorID := vars["doctor_id"]
	date := vars["date"]

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	id, err := strconv.Atoi(doctorID)
	if err != nil {
		http.Error(w, "Неверный ID врача", http.StatusBadRequest)
		return
	}

	var slots []models.TimeSlot
	result := database.DB.
		Where("doctor_id = ? AND slot_date = ? AND is_available = ? AND appointment_id IS NULL", id, date, true).
		Order("start_time asc").
		Find(&slots)

	if result.Error != nil {
		http.Error(w, "Ошибка получения слотов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(slots)
}

// CreateAppointment - создать запись
func CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var req CreateAppointmentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Проверяем слот
	var timeSlot models.TimeSlot
	if err := database.DB.First(&timeSlot, req.TimeSlotID).Error; err != nil {
		http.Error(w, "Слот не найден", http.StatusNotFound)
		return
	}

	if !timeSlot.IsAvailable || timeSlot.AppointmentID != nil {
		http.Error(w, "Слот уже занят", http.StatusConflict)
		return
	}

	// Начинаем транзакцию
	tx := database.DB.Begin()

	// 1. Находим или создаём клиента
	var client models.Client
	var clientID uint

	if req.ClientID != nil && *req.ClientID > 0 {
		// Используем существующего клиента
		if err := tx.First(&client, *req.ClientID).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Клиент не найден", http.StatusNotFound)
			return
		}
		clientID = client.ID
	} else {
		// Проверяем, есть ли клиент с таким телефоном
		err := tx.Where("phone = ?", req.Phone).First(&client).Error
		if err == nil {
			// Клиент существует
			clientID = client.ID
		} else {
			// Создаём нового клиента
			newClient := models.Client{
				Phone:           req.Phone,
				LastName:        req.LastName,
				FirstName:       req.FirstName,
				MiddleName:      req.MiddleName,
				BirthDate:       req.BirthDate,
				ContactPhone:    req.ContactPhone,
				ContactTelegram: req.ContactTelegram,
				ContactWhatsApp: req.ContactWhatsApp,
			}

			if err := tx.Create(&newClient).Error; err != nil {
				tx.Rollback()
				http.Error(w, "Ошибка создания клиента", http.StatusInternalServerError)
				return
			}
			clientID = newClient.ID
		}
	}

	// 2. Создаём запись на приём
	appointment := models.Appointment{
		ClientID:   clientID,
		DoctorID:   req.DoctorID,
		TimeSlotID: req.TimeSlotID,
		Status:     "pending",
	}

	if err := tx.Create(&appointment).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Ошибка создания записи", http.StatusInternalServerError)
		return
	}

	// 3. Обновляем слот
	if err := tx.Model(&timeSlot).Updates(map[string]interface{}{
		"is_available":   false,
		"appointment_id": appointment.ID,
	}).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Ошибка обновления слота", http.StatusInternalServerError)
		return
	}

	tx.Commit()

	// Загружаем полную информацию для ответа
	var result models.Appointment
	database.DB.
		Preload("Client").
		Preload("Doctor").
		Preload("Doctor.Specialization").
		Preload("TimeSlot").
		First(&result, appointment.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

// GetAppointment - получить информацию о записи
func GetAppointment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var appointment models.Appointment
	result := database.DB.
		Preload("Client").
		Preload("Doctor").
		Preload("Doctor.Specialization").
		Preload("TimeSlot").
		First(&appointment, id)

	if result.Error != nil {
		http.Error(w, "Запись не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointment)
}

// GetClientAppointments - получить все записи клиента
func GetClientAppointments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID := vars["client_id"]

	var appointments []models.Appointment
	result := database.DB.
		Where("client_id = ?", clientID).
		Preload("Doctor").
		Preload("Doctor.Specialization").
		Preload("TimeSlot").
		Order("created_at desc").
		Find(&appointments)

	if result.Error != nil {
		http.Error(w, "Ошибка получения записей", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointments)
}

// CancelAppointment - отменить запись
func CancelAppointment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	tx := database.DB.Begin()

	var appointment models.Appointment
	if err := tx.First(&appointment, id).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Запись не найдена", http.StatusNotFound)
		return
	}

	// Обновляем статус записи
	if err := tx.Model(&appointment).Update("status", "cancelled").Error; err != nil {
		tx.Rollback()
		http.Error(w, "Ошибка отмены", http.StatusInternalServerError)
		return
	}

	// Освобождаем слот
	if err := tx.Model(&models.TimeSlot{}).Where("id = ?", appointment.TimeSlotID).
		Updates(map[string]interface{}{
			"is_available":   true,
			"appointment_id": nil,
		}).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Ошибка освобождения слота", http.StatusInternalServerError)
		return
	}

	tx.Commit()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Запись отменена"})
}

// GetAvailableDates - получить доступные даты для врача
func GetAvailableDates(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	doctorID := vars["doctor_id"]

	id, err := strconv.Atoi(doctorID)
	if err != nil {
		http.Error(w, "Неверный ID врача", http.StatusBadRequest)
		return
	}

	var dates []string
	// Получаем уникальные даты, где есть свободные слоты
	result := database.DB.Model(&models.TimeSlot{}).
		Where("doctor_id = ? AND is_available = ? AND slot_date >= ?",
			id, true, time.Now().Format("2006-01-02")).
		Group("slot_date").
		Order("slot_date").
		Pluck("slot_date", &dates)

	if result.Error != nil {
		http.Error(w, "Ошибка получения дат", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(dates)
}
