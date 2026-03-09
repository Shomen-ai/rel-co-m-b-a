// internal/handlers/admin_handlers.go
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"healthcare/internal/database"
	"healthcare/internal/models"
)

func GetAdminStats(w http.ResponseWriter, r *http.Request) {
	type Stats struct {
		TodayAppointments int64 `json:"todayAppointments"`
		ActiveDoctors     int64 `json:"activeDoctors"`
		TotalClients      int64 `json:"totalClients"`
	}

	var stats Stats
	today := time.Now().Format("2006-01-02")

	// Записи на сегодня
	database.DB.Model(&models.Appointment{}).
		Where("appointment_date = ?", today).
		Count(&stats.TodayAppointments)

	// Активные врачи
	database.DB.Model(&models.Doctor{}).
		Where("is_active = ?", true).
		Count(&stats.ActiveDoctors)

	// Всего клиентов
	database.DB.Model(&models.Client{}).
		Count(&stats.TotalClients)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func GetAdminAppointments(w http.ResponseWriter, r *http.Request) {
	var appointments []struct {
		ID              uint   `json:"id"`
		LastName        string `json:"last_name"`
		FirstName       string `json:"first_name"`
		DoctorName      string `json:"doctor_name"`
		AppointmentDate string `json:"appointment_date"`
		AppointmentTime string `json:"appointment_time"`
		Status          string `json:"status"`
	}

	// Проверяем соединение с БД
	if database.DB == nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	// Выполняем запрос с отладкой
	result := database.DB.Table("appointments").
		Select(`appointments.id, 
                clients.last_name, 
                clients.first_name,
                CONCAT(doctors.last_name, ' ', doctors.first_name) as doctor_name,
                appointments.appointment_date,
                appointments.appointment_time,
                appointments.status`).
		Joins("JOIN clients ON appointments.client_id = clients.id").
		Joins("JOIN doctors ON appointments.doctor_id = doctors.id").
		Order("appointments.created_at DESC").
		Limit(10).
		Scan(&appointments)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Если записей нет, возвращаем пустой массив
	if appointments == nil {
		appointments = []struct {
			ID              uint   `json:"id"`
			LastName        string `json:"last_name"`
			FirstName       string `json:"first_name"`
			DoctorName      string `json:"doctor_name"`
			AppointmentDate string `json:"appointment_date"`
			AppointmentTime string `json:"appointment_time"`
			Status          string `json:"status"`
		}{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointments)
}
