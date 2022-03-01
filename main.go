package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
	"github.com/mahjadan/go-mfa/handlers"
	"github.com/mahjadan/go-mfa/pkg/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var database *mongo.Database

type goKey string

const sessUserKey = "GO-MFA-USERNAME"

//docker run -p 27017:27017 mongo:4.4
func main() {
	database = initMongoDB().Database("mfa-demo")
	repo := repository.NewMongo(database)
	handler := handlers.New(repo, session.New(), sessUserKey)

	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/statics", "./public")
	app.Use(recover.New())
	app.Use(logger.New())

	app.Get("/", handler.HandleDashboard)
	app.Get("/dashboard", handler.HandleDashboard)
	app.Post("/enable-mfa", handler.HandleEnableMFA)
	app.Get("/profile", handler.HandleProfile)
	app.Post("/register", handler.HandleRegister)
	app.Get("/register", handler.RegisterPage)
	app.Post("/login", handler.HandleLogin)
	app.Post("/mfa", handler.HandleMFA)
	app.Get("/login", handler.LoginPage)
	app.Get("/logout", handler.HandleLogout)

	app.Listen(os.Getenv("PORT"))

}

func initMongoDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
