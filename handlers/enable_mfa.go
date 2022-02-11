package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mahjadan/login/cmd/handle"
)

func (h Handler) HandleEnableMFA(c *fiber.Ctx) error {
	sess, err := h.session.Get(c)
	if err != nil {
		return c.JSON(handle.NewInternalServerResponse(err.Error()))
	}
	username := sess.Get(h.sessUserKey)
	return c.Render("mfa", fiber.Map{
		"mfa":      "active",
		"username": username,
	})
}

//todo create mfa page to show the mfa options for the user
