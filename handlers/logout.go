package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mahjadan/login/cmd/handle"
)

func (h Handler) HandleLogout(c *fiber.Ctx) error {
	sess, err := h.session.Get(c)
	if err != nil {
		return c.JSON(handle.NewInternalServerResponse(err.Error()))
	}
	sess.Destroy()
	return c.Redirect("/")
}
