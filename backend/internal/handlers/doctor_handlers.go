package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"healthcare/internal/database"
	"healthcare/internal/models"
)

// GetAllDoctors возвращает всех активных врачей
func GetAllDoctors(w http.ResponseWriter, r *http.Request) {
	var doctors []models.Doctor

	result := database.DB.Where("is_active = ?", true).
		Order("display_order asc").
		Find(&doctors)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doctors)
}

func GetDoctorByID(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")

	var idStr string
	for i, part := range pathParts {
		if part == "doctors" && i+1 < len(pathParts) {
			idStr = pathParts[i+1]
			break
		}
	}

	if idStr == "" {
		http.Error(w, "ID врача не указан", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID врача", http.StatusBadRequest)
		return
	}

	var doctor models.Doctor

	result := database.DB.
		Preload("Specialization").
		Preload("Services").
		First(&doctor, id)

	if result.Error != nil {
		http.Error(w, "Врач не найден", http.StatusNotFound)
		return
	}

	if !doctor.IsActive {
		http.Error(w, "Врач не активен", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doctor)
}

// GetDoctorsBySpecialization возвращает врачей по ID специализации
func GetDoctorsBySpecialization(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из URL: /api/doctors/by-specialization/7
	pathParts := strings.Split(r.URL.Path, "/")

	// Проверяем, что путь правильный
	if len(pathParts) < 5 {
		http.Error(w, "ID специализации не указан", http.StatusBadRequest)
		return
	}

	// Берем последнюю часть пути - это ID специализации
	specIDStr := pathParts[len(pathParts)-1]

	// Преобразуем строку в число
	specID, err := strconv.Atoi(specIDStr)
	if err != nil || specID <= 0 {
		http.Error(w, "Неверный ID специализации", http.StatusBadRequest)
		return
	}

	// Ищем врачей в базе данных
	var doctors []models.Doctor

	result := database.DB.
		Where("specialization_id = ? AND is_active = ?", specID, true).
		Preload("Specialization").
		Order("display_order asc, last_name asc").
		Find(&doctors)

	if result.Error != nil {
		http.Error(w, "Ошибка при получении врачей", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doctors)
}
