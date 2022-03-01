package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	mo "github.com/mahjadan/go-mfa/pkg/models"
	"github.com/mahjadan/login/cmd/handle"
)

func (h Handler) RegisterPage(c *fiber.Ctx) error {
	return c.Render("register", nil)
}

func (h Handler) HandleRegister(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	fmt.Println(req)
	exists, err := h.repo.Exists(c.Context(), req.Username)
	if exists && err == nil {
		fmt.Println(err)
		return c.Render("register", fiber.Map{"error": "user already exists"})
	}
	mClient := mo.New(req.Username, req.Password)
	err = h.repo.Set(c.Context(), &mClient)
	if err != nil {
		return c.JSON(handle.NewInternalServerResponse(err.Error()))
	}
	return c.Redirect("/login")
}
