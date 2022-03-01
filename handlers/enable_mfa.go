package handlers

import (
	"github.com/gofiber/fiber/v2"
	mo "github.com/mahjadan/go-mfa/pkg/models"
	"github.com/mahjadan/go-mfa/pkg/service"
	"github.com/mahjadan/login/cmd/handle"
	"time"
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
	var mClient mo.MongoClient
	mClient.Username = username.(string)
	mfa := mo.Mfa{
		Type:           "App Authentication",
		DisplayName:    "RD Station Accounts",
		Data:           map[string]string{"secret": secret},
		ActivationTime: time.Now(),
	}
	mClient.MFA = append(mClient.MFA, mfa)
	if value == "on" {
		err := h.repo.Update(c.Context(), mClient)
		if err != nil {
			return c.JSON(handle.NewInternalServerResponse(err.Error()))
		}
	}
	return c.Redirect("/profile")
}
