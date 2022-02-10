package handlers

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/mahjadan/go-mfa/pkg/repository"
)

func New(repo *repository.MongoRepo, store *session.Store, sessionKey string) Handler {
	return Handler{
		repo:        repo,
		session:     store,
		sessUserKey: sessionKey,
	}
}

type Handler struct {
	repo        *repository.MongoRepo
	session     *session.Store
	sessUserKey string
}
