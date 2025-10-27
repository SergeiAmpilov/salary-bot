package handler

import (
	"salary-bot/internal/user/repository"

	"github.com/gofiber/fiber/v2"
)

type statHandler struct {
	userRepo repository.Repository
}

func NewStatHandler(userRepo repository.Repository) Handler {
	return &statHandler{userRepo: userRepo}
}

func (h *statHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.userRepo.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch users"})
	}

	return c.JSON(users)
}

func (h *statHandler) GetStats(c *fiber.Ctx) error {
	last24h, last7d, err := h.userRepo.GetNewUsersStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch stats",
		})
	}

	return c.JSON(fiber.Map{
		"new_users": fiber.Map{
			"last_24_hours": last24h,
			"last_7_days":   last7d,
		},
	})
}
