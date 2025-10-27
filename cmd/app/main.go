// cmd/app/main.go
package main

import (
	"log"
	"salary-bot/internal/bot"
	"salary-bot/internal/config"
	"salary-bot/internal/router"
	"salary-bot/internal/salary/handler"
	"salary-bot/internal/salary/repository"
	"salary-bot/internal/salary/service"
	stathandler "salary-bot/internal/stat/handler"
	"salary-bot/internal/storage"
	userrepository "salary-bot/internal/user/repository"

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

	userRepo := userrepository.New(db.DB)

	// Слои
	sRepo := repository.New(db.DB)
	sSvc := service.New(sRepo)
	sHandler := handler.NewSalaryHandler(sSvc)
	statHandler := stathandler.NewStatHandler(userRepo)

	// Telegram Bot
	tgBot, err := bot.NewBot(cfg.TelegramBotToken, sSvc, userRepo)
	if err != nil {
		log.Fatal("Не удалось инициализировать Telegram бота:", err)
	}

	// Запуск бота в отдельной горутине
	go tgBot.Start()

	// Fiber
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	// Роуты
	router.SetupRoutes(app, sHandler, statHandler)

	// Запуск
	log.Printf("Сервер запущен на порту %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
