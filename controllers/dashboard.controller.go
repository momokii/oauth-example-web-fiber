package controllers

import "github.com/gofiber/fiber/v2"

func DashboardView(c *fiber.Ctx) error {
	username := c.Locals("username").(string)

	return c.Render("index", fiber.Map{
		"Error": nil,
		"User":  username,
	})
}
