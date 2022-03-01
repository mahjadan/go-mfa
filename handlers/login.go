package handlers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	mo "github.com/mahjadan/go-mfa/pkg/models"
	"github.com/mahjadan/login/cmd/handle"
	"github.com/xlzd/gotp"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
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
		Code     string
	}{}
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	ctx, cancelFunc := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancelFunc()
	fmt.Println(req)

	var mClient mo.MongoClient
	if strings.TrimSpace(req.Code) == "" {
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
			return c.Render("login", fiber.Map{
				"withOPT": true,
			})
		}
		// save login session
		sess.Set(h.sessUserKey, req.Username)
		defer sess.Save()
		return c.Redirect("/dashboard")
	}

	// verify the code and create session and redirect to /dashboard
	fmt.Println(req.Code)
	sess, err := h.session.Get(c)
	if err != nil {
		panic(err)
	}
	username := sess.Get("mfa-session")
	err = h.repo.Get(ctx, username.(string), &mClient)
	if err != nil {
		fmt.Println(err)
		return c.JSON(handle.NewInternalServerResponse(err.Error()))
	}
	totp := gotp.NewDefaultTOTP(mClient.MFA[0].Data["secret"])
	verified := totp.Verify(req.Code, int(time.Now().Unix()))
	fmt.Println("verified: ", verified)
	if verified {
		sess.Set(h.sessUserKey, req.Username)
		defer sess.Save()
		return c.Redirect("/dashboard")
	}
	return c.Redirect("/login")
}
