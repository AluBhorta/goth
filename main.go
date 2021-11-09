package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app := fiber.New()

	// routes
	app.Get("/", index)

	listenHost := os.Getenv("LISTEN_ON_HOST")
	listenPort := os.Getenv("LISTEN_ON_PORT")

	if err := app.Listen(listenHost + ":" + listenPort); err != nil {
		log.Fatalln(err)
	}
}

func index(c *fiber.Ctx) error {
	log.Println("serving index")
	return c.JSON(fiber.Map{"message": "API is functional 🚀"})
}
