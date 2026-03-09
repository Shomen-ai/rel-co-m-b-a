package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"healthcare/internal/database"
	"healthcare/internal/handlers"
)

func main() {
	// Подключаемся к базе данных
	database.Connect()

	// Создаём роутер
	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok", "message": "Server is running"}`))
	}).Methods("GET")

	// Специализации
	r.HandleFunc("/api/specializations", handlers.GetAllSpecializations).Methods("GET")
	r.HandleFunc("/api/specializations/{id}", handlers.GetSpecializationByID).Methods("GET")
	r.HandleFunc("/api/specializations/{id}/services", handlers.GetServicesBySpecialization).Methods("GET")

	// Врачи
	r.HandleFunc("/api/doctors", handlers.GetAllDoctors).Methods("GET")
	r.HandleFunc("/api/doctors/{id}", handlers.GetDoctorByID).Methods("GET")
	r.HandleFunc("/api/doctors/by-specialization/{id}", handlers.GetDoctorsBySpecialization).Methods("GET")

	// Услуги
	r.HandleFunc("/api/services", handlers.GetAllServices).Methods("GET")
	r.HandleFunc("/api/services/{id}", handlers.GetServiceByID).Methods("GET")

	// Расписание (слоты)
	r.HandleFunc("/api/slots/{date}", handlers.GetAvailableSlots).Methods("GET")
	r.HandleFunc("/api/doctors/{doctor_id}/slots/{date}", handlers.GetDoctorSlots).Methods("GET")

	// Записи на приём
	r.HandleFunc("/api/appointments", handlers.CreateAppointment).Methods("POST")
	r.HandleFunc("/api/appointments/{id}", handlers.GetAppointment).Methods("GET")
	r.HandleFunc("/api/appointments/{id}/cancel", handlers.CancelAppointment).Methods("PUT")
	r.HandleFunc("/api/clients/{client_id}/appointments", handlers.GetClientAppointments).Methods("GET")

	r.HandleFunc("/api/doctors/{doctor_id}/available-dates", handlers.GetAvailableDates).Methods("GET")

	// Маршруты для админки
	r.HandleFunc("/api/admin/stats", handlers.GetAdminStats).Methods("GET")
	r.HandleFunc("/api/admin/appointments", handlers.GetAdminAppointments).Methods("GET")
	// Запускаем сервер
	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
