package handlers

import (
	"github.com/gofiber/fiber/v2"
	mo "github.com/mahjadan/go-mfa/pkg/models"
	"github.com/mahjadan/go-mfa/pkg/service"
	"github.com/mahjadan/login/cmd/handle"
)

func (h Handler) HandleEnableMFA(c *fiber.Ctx) error {
	value := c.FormValue("enable-auth")
	sess, err := h.session.Get(c)
	if err != nil {
		return c.JSON(handle.NewInternalServerResponse(err.Error()))
	}
	username := sess.Get(h.sessUserKey)
	if username == nil {
		return c.Redirect("/login")
	}
	secret := service.GenerateOtpSecret()

	mClient := mo.NewWithMFA(username.(string), secret)
	if value == "on" {
		err := h.repo.Update(c.Context(), mClient)
		if err != nil {
			return c.JSON(handle.NewInternalServerResponse(err.Error()))
		}
	}
	return c.Redirect("/profile")
}
