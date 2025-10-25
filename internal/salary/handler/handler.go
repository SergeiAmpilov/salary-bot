// internal/salary/handler/handler.go
package handler

import (
	"salary-bot/internal/salary/model"
	"salary-bot/internal/salary/service"

	"github.com/gofiber/fiber/v2"
)

type salaryHandler struct {
	service service.Service
}

func NewSalaryHandler(svc service.Service) Handler {
	return &salaryHandler{service: svc}
}

func (h *salaryHandler) Add(c *fiber.Ctx) error {
	var dto model.CreateSalaryDTO

	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON",
		})
	}

	if err := h.service.Create(&dto); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save salary data",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Salary record created successfully",
	})
}

func (h *salaryHandler) List(c *fiber.Ctx) error {
	salaries, err := h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch salary data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(salaries)
}
