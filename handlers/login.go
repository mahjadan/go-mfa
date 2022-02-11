package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	mo "github.com/mahjadan/go-mfa/pkg/models"
	"github.com/mahjadan/login/cmd/handle"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h Handler) LoginPage(c *fiber.Ctx) error {
	err := c.Query("error")
	return c.Render("login", fiber.Map{"error": err})
}

func (h Handler) HandleLogin(c *fiber.Ctx) error {
	req := struct {
		Username string
	}{}
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	fmt.Println(req)
	var mClient mo.MongoClient
	err := h.repo.Get(c.Context(), req.Username, &mClient)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Render("login", fiber.Map{
				"error": "user not found",
			})
		}
		fmt.Println(err)
		return c.JSON(handle.NewInternalServerResponse(err.Error()))
	}
	// save session
	sess, err := h.session.Get(c)
	if err != nil {
		panic(err)
	}
	sess.Set(h.sessUserKey, req.Username)
	defer sess.Save()
	return c.Redirect("/dashboard")
}
