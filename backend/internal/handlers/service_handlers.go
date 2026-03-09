package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"healthcare/internal/database"
	"healthcare/internal/models"
)

// GetAllServices возвращает все активные услуги
func GetAllServices(w http.ResponseWriter, r *http.Request) {
	var services []models.Service

	// Базовый запрос
	query := database.DB.Where("is_active = ?", true).
		Preload("Specialization"). // Загружаем специализацию вместо Doctor
		Order("display_order asc")

	// Фильтр по специализации, если передан параметр specialization_id
	specializationID := r.URL.Query().Get("specialization_id")
	if specializationID != "" {
		id, err := strconv.Atoi(specializationID)
		if err == nil && id > 0 {
			query = query.Where("specialization_id = ?", id)
		}
	}

	// Фильтр по поиску
	search := r.URL.Query().Get("search")
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	result := query.Find(&services)

	if result.Error != nil {
		http.Error(w, "Ошибка при получении услуг", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

// GetServiceByID возвращает услугу по ID
func GetServiceByID(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из URL: /api/services/1
	pathParts := strings.Split(r.URL.Path, "/")

	if len(pathParts) < 3 {
		http.Error(w, "ID услуги не указан", http.StatusBadRequest)
		return
	}

	idStr := pathParts[len(pathParts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID услуги", http.StatusBadRequest)
		return
	}

	var service models.Service
	result := database.DB.
		Preload("Specialization").
		First(&service, id)

	if result.Error != nil {
		http.Error(w, "Услуга не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(service)
}

// GetServicesBySpecialization возвращает услуги по ID специализации
func GetServicesBySpecialization(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из URL: /api/specializations/7/services
	pathParts := strings.Split(r.URL.Path, "/")

	if len(pathParts) < 4 {
		http.Error(w, "ID специализации не указан", http.StatusBadRequest)
		return
	}

	specIDStr := pathParts[len(pathParts)-2] // Берем предпоследнюю часть (перед services)
	specID, err := strconv.Atoi(specIDStr)
	if err != nil || specID <= 0 {
		http.Error(w, "Неверный ID специализации", http.StatusBadRequest)
		return
	}

	var services []models.Service
	result := database.DB.
		Where("specialization_id = ? AND is_active = ?", specID, true).
		Preload("Specialization").
		Order("display_order asc, name asc").
		Find(&services)

	if result.Error != nil {
		http.Error(w, "Ошибка при получении услуг", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}
