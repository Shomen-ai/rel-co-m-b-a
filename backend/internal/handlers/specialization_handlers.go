package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"healthcare/internal/database"
	"healthcare/internal/models"

	"github.com/gorilla/mux"
)

// GetAllSpecializations возвращает все специализации
func GetAllSpecializations(w http.ResponseWriter, r *http.Request) {
	var specializations []models.Specialization

	// Получаем все специализации из базы, сортируем по имени
	result := database.DB.Order("name asc").Find(&specializations)

	if result.Error != nil {
		http.Error(w, "Ошибка при получении специализаций", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(specializations)
}

func GetSpecializationByID(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из URL (через mux)
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID специализации", http.StatusBadRequest)
		return
	}

	var specialization models.Specialization
	result := database.DB.First(&specialization, id)

	if result.Error != nil {
		http.Error(w, "Специализация не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(specialization)
}
