package handlers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	mo "github.com/mahjadan/go-mfa/pkg/models"
	"github.com/mahjadan/login/cmd/handle"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (h Handler) LoginPage(c *fiber.Ctx) error {
	error := c.Query("error")
	return c.Render("login", fiber.Map{"error": error})
}

func (h Handler) HandleLogin(c *fiber.Ctx) error {
	req := struct {
		Username string
		Password string
	}{}
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	ctx, cancelFunc := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancelFunc()
	fmt.Println(req)

	var mClient mo.MongoClient
	// normal login, not checking for password for simplicity
	err := h.repo.Get(ctx, req.Username, &mClient)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Render("login", fiber.Map{
				"error": "user not found",
			})
		}
		fmt.Println(err)
		return c.JSON(handle.NewInternalServerResponse(err.Error()))
	}
	sess, err := h.session.Get(c)
	if err != nil {
		panic(err)
	}

	if len(mClient.MFA) != 0 {
		sess.Set("mfa-session", req.Username)
		defer sess.Save()
		return c.Render("mfa", nil)
	}
	// save login session
	sess.Set(h.sessUserKey, req.Username)
	defer sess.Save()
	return c.Redirect("/dashboard")
}
