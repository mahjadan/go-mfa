package handlers

import "github.com/gofiber/fiber/v2"

func (h Handler) HandleHome(c *fiber.Ctx) error {
	return c.Render("home", nil)
}
