// internal/router/router.go
package router

import (
	"salary-bot/internal/salary/handler"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes регистрирует все маршруты приложения
func SetupRoutes(app *fiber.App, sHandler handler.Handler) {

	// Тестовый эндпоинт
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "OK",
		})
	})

	// Эндпоинт для добавления данных о вакансии
	app.Post("/vacancy", sHandler.Add)
	app.Get("/salary", sHandler.List)
	app.Post("/filter", sHandler.Filter)
}
