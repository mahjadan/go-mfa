package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h Handler) HandleDashboard(c *fiber.Ctx) error {
	sess, err := h.session.Get(c)
	if err != nil {
		panic(err)
	}
	username := sess.Get(h.sessUserKey)
	if username == nil {
		return c.Redirect("/login")
	}

	return c.Render("dashboard", fiber.Map{
		"username":      username,
		"dashboardMenu": "active",
	})
}
