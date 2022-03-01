package handlers

import (
	"github.com/gofiber/fiber/v2"
	mo "github.com/mahjadan/go-mfa/pkg/models"
	"github.com/mahjadan/go-mfa/pkg/service"
	"github.com/mahjadan/login/cmd/handle"
)

func (h Handler) HandleProfile(c *fiber.Ctx) error {
	sess, err := h.session.Get(c)
	if err != nil {
		return c.JSON(handle.NewInternalServerResponse(err.Error()))
	}
	username := sess.Get(h.sessUserKey)
	if username == nil {
		return c.Redirect("/login")
	}
	var mClient mo.MongoClient
	err = h.repo.Get(c.Context(), username.(string), &mClient)
	if err != nil {
		return c.Render("profile", fiber.Map{
			"profileMenu": "active",
			"error":       err})
	}
	var mfaEnabled bool
	var mfa mo.Mfa
	if len(mClient.MFA) != 0 {
		mfaEnabled = true
		mfa = mClient.MFA[0]
	}

	if !mfaEnabled {
		return c.Render("profile", fiber.Map{
			"show_enable_authentication": true,
			"profileMenu":                "active",
		})
	}
	url := service.GenerateURL(mfa.Data["secret"], username.(string))
	return c.Render("profile", fiber.Map{
		"profileMenu":                "active",
		"show_enable_authentication": false,
		"url":                        url,
	})
}

//todo create mfa page to show the mfa options for the user
