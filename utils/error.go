package utils

import "github.com/gofiber/fiber/v2"

func ErrorJSON(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error":   true,
		"message": message,
	})
}
