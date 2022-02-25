package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
	"github.com/mahjadan/go-mfa/handlers"
	"github.com/mahjadan/go-mfa/pkg/models"
	"github.com/mahjadan/go-mfa/pkg/repository"
	"github.com/mahjadan/login/cmd/handle"
	"github.com/xlzd/gotp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var secretLength = 16

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
	app.Post("/auth-mfa", handleAuth)
	app.Get("/enable-mfa", handler.HandleEnableMFA)
	app.Get("/users/:email", handleGet)
	app.Post("/register", handler.HandleRegister)
	app.Get("/register", handler.RegisterPage)
	app.Post("/login", handler.HandleLogin)
	app.Get("/login", handler.LoginPage)
	app.Get("/logout", handler.HandleLogout)

	secret := gotp.RandomSecret(secretLength)

	fmt.Println("Current OTP is", gotp.NewDefaultTOTP(secret).Now())

	app.Listen(os.Getenv("PORT"))

}

func handleGet(c *fiber.Ctx) error {
	email := c.Params("email")
	var mClient mo.MongoClient
	err := database.Collection("clients").FindOne(c.Context(), bson.M{"username": email}).Decode(&mClient)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(handle.NotFoundErrorResponse)
		}
		fmt.Println(err)
		return c.JSON(handle.NewInternalServerResponse(err.Error()))
	}
	return c.JSON(mClient)
}

func handleAuth(c *fiber.Ctx) error {
	var mfaReq mo.AuthRequest
	if err := c.BodyParser(&mfaReq); err != nil {
		return err
	}
	fmt.Println(mfaReq)
	var mClient mo.MongoClient
	err := database.Collection("clients").FindOne(c.Context(), bson.M{"username": mfaReq.Username}).Decode(&mClient)
	if err != nil {
		fmt.Println(err)
		return c.JSON(handle.NotFoundErrorResponse)
	}
	return c.JSON(mClient)
}

func initMongoDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
