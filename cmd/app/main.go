// cmd/app/main.go
package main

import (
	"log"
	"salary-bot/internal/config"
	"salary-bot/internal/router"
	"salary-bot/internal/salary/handler"
	"salary-bot/internal/salary/repository"
	"salary-bot/internal/salary/service"
	"salary-bot/internal/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Конфиг
	cfgReader := config.NewEnvConfigService()
	cfg, err := cfgReader.Read()
	if err != nil {
		log.Fatal("Не удалось загрузить конфигурацию:", err)
	}

	// БД
	db := storage.NewStorage("./data/salary.db")
	defer db.DB.Close()

	// Слои
	sRepo := repository.New(db.DB)
	sSvc := service.New(sRepo)
	sHandler := handler.NewSalaryHandler(sSvc)

	// Fiber
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	// Роуты
	router.SetupRoutes(app, sHandler)

	// Запуск
	log.Printf("Сервер запущен на порту %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
