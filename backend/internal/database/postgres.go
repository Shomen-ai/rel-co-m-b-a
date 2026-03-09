package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"healthcare/internal/models"
)

var DB *gorm.DB

// Connect инициализирует подключение к базе данных через GORM
func Connect() {
	// Загружаем .env файл
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Файл .env не найден, использую системные переменные окружения")
	}

	// Формируем DSN из отдельных переменных (более гибко, чем DB_URL)
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Проверяем, что все переменные заданы
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("❌ Не все переменные окружения для БД заданы. Проверьте .env файл")
	}

	// Формируем строку подключения
	dsn := "host=" + host + " port=" + port + " user=" + user +
		" password=" + password + " dbname=" + dbname + " sslmode=disable"

	// Настраиваем подключение
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Логируем SQL запросы
	})
	if err != nil {
		log.Fatal("❌ Ошибка подключения к БД через GORM:", err)
	}

	// Автоматическая миграция (создание таблиц по моделям)
	err = DB.AutoMigrate(
		&models.Doctor{},
		&models.Service{},
		&models.Client{},
		&models.Appointment{},
	)
	if err != nil {
		log.Fatal("❌ Ошибка миграции:", err)
	}

	log.Println("✅ Успешное подключение к базе данных PostgreSQL и миграция выполнена!")
}
