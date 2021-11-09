package main

import (
	"log"
	"os"

	authapi "github.com/alubhorta/goth/api/auth"
	userapi "github.com/alubhorta/goth/api/user"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app := fiber.New()

	setupRoutes(app)

	listenHost := os.Getenv("LISTEN_ON_HOST")
	listenPort := os.Getenv("LISTEN_ON_PORT")

	if err := app.Listen(listenHost + ":" + listenPort); err != nil {
		log.Fatalln(err)
	}
}

func setupRoutes(app *fiber.App) {
	app.Get("/", index)

	// auth routes
	app.Post("/api/v1/auth/signup", authapi.Signup)
	app.Post("/api/v1/auth/login", authapi.Login)
	app.Post("/api/v1/auth/logout", authapi.Logout)
	app.Post("/api/v1/auth/refresh", authapi.Refresh)
	app.Post("/api/v1/auth/reset/init", authapi.ResetPasswordInit)
	app.Post("/api/v1/auth/reset/verify", authapi.ResetPasswordVerify)
	app.Delete("/api/v1/auth/delete/:id", authapi.DeleteAccount)

	// user routes
	app.Get("/api/v1/user", userapi.GetAll)
	app.Get("/api/v1/user/:id", userapi.GetOne)
	app.Put("/api/v1/user/:id", userapi.UpdateOne)

	// TODO: add requiresAuth middleware
}

func index(c *fiber.Ctx) error {
	log.Println("serving index...")
	return c.JSON(fiber.Map{"message": "API is functional 🚀"})
}
