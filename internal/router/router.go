// internal/router/router.go
package router

import (
	"salary-bot/internal/salary/handler"
	shandler "salary-bot/internal/stat/handler"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes регистрирует все маршруты приложения
func SetupRoutes(
	app *fiber.App,
	sHandler handler.Handler,
	statHandler shandler.Handler,
) {

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

	// stat
	app.Get("/users", statHandler.GetUsers)
	app.Get("/stat", statHandler.GetStats)
}
