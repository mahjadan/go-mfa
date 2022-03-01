package handlers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	mo "github.com/mahjadan/go-mfa/pkg/models"
	"github.com/mahjadan/go-mfa/pkg/service"
	"github.com/mahjadan/login/cmd/handle"
	"time"
)

func (h Handler) HandleMFA(c *fiber.Ctx) error {
	req := struct {
		Code string
	}{}
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	ctx, cancelFunc := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancelFunc()

	fmt.Println(req.Code)
	sess, err := h.session.Get(c)
	if err != nil {
		panic(err)
	}
	username := sess.Get("mfa-session")
	var mClient mo.MongoClient
	err = h.repo.Get(ctx, username.(string), &mClient)
	if err != nil {
		fmt.Println(err)
		return c.JSON(handle.NewInternalServerResponse(err.Error()))
	}
	// we know for sure that he has MFA object, otherwise he will not have mfa-session
	verified := service.VerifyOtpCode(req.Code, mClient.MFA[0].Data["secret"])
	if verified {
		sess.Set(h.sessUserKey, username)
		sess.Delete("mfa-session")
		defer sess.Save()
		return c.Redirect("/dashboard")
	}
	return c.Render("mfa", fiber.Map{
		"error": "invalid code",
	})
}
