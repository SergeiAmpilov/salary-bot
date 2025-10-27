package handler

import "github.com/gofiber/fiber/v2"

type Handler interface {
	GetUsers(c *fiber.Ctx) error
	GetStats(c *fiber.Ctx) error
}
