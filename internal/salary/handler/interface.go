// internal/vacancy/handler/interface.go
package handler

import "github.com/gofiber/fiber/v2"

// Handler описывает контракт для обработчиков вакансий
type Handler interface {
	Add(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
	Filter(c *fiber.Ctx) error
}
