package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h Handler) HandleDashboard(c *fiber.Ctx) error {
	username := c.Query("username")
	sess, err := h.session.Get(c)
	if err != nil {
		panic(err)
	}
	if name := sess.Get(h.sessUserKey); name != username {
		return c.Redirect("/login")
	}

	return c.Render("dashboard", fiber.Map{"username": username})
}
