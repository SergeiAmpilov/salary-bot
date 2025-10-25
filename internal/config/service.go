// internal/config/service.go
package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config хранит все переменные окружения приложения
type Config struct {
	Port string
}

// envConfigService реализует ConfigReader
type envConfigService struct{}

// NewEnvConfigService создаёт новый экземпляр сервиса конфигурации
func NewEnvConfigService() ConfigReader {
	return &envConfigService{}
}

// Read загружает .env и возвращает структуру Config
func (s *envConfigService) Read() (*Config, error) {
	// Загружаем .env файл (игнорируем ошибку, если файла нет — допустимо в продакшене)
	_ = godotenv.Load()

	return &Config{
		Port: os.Getenv("PORT"),
	}, nil
}
