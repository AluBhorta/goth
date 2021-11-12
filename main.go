package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	authapi "github.com/alubhorta/goth/api/auth"
	userapi "github.com/alubhorta/goth/api/user"
	commonclients "github.com/alubhorta/goth/models/common"
	tokenutils "github.com/alubhorta/goth/utils/token"

	"github.com/alubhorta/goth/db/cacheclient"
	"github.com/alubhorta/goth/db/dbclient"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app := fiber.New()
	// TODO: add cors

	dbclient := &dbclient.MongoDbClient{}
	dbclient.Init()

	redisClient := &cacheclient.RedisClient{}
	redisClient.Init()

	userCtx := context.WithValue(
		context.Background(),
		commonclients.CommonClients{},
		&commonclients.CommonClients{
			DbClient:    dbclient,
			CacheClient: redisClient,
		},
	)
	app.Use(func(c *fiber.Ctx) error {
		c.SetUserContext(userCtx)
		return c.Next()
	})

	setupRoutes(app)

	// ensure cleanup
	cleanupFunc := func() {
		log.Println("running cleanup tasks...")
		redisClient.Cleanup()
		dbclient.Cleanup(userCtx)
		app.Shutdown()
		log.Println("all done! bye 👋")
	}
	ensureGracefulTermination(cleanupFunc)

	// start serving!
	listenHost := os.Getenv("GOTH_LISTEN_HOST")
	listenPort := os.Getenv("GOTH_LISTEN_PORT")

	if err := app.Listen(listenHost + ":" + listenPort); err != nil {
		cleanupFunc()
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
	app.Delete("/api/v1/auth/delete", authapi.DeleteAccount)

	// user routes
	app.Get("/api/v1/user/:id", tokenutils.RequiresAuth, userapi.GetOne)
	app.Put("/api/v1/user/:id", tokenutils.RequiresAuth, userapi.UpdateOne)
	// TODO: remove id from path and use id parsed from access token
}

func index(c *fiber.Ctx) error {
	log.Println("serving index...")
	return c.JSON(fiber.Map{"message": "API is functional 🚀"})
}

func ensureGracefulTermination(cleanupFunc func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		s := <-c
		log.Printf("gracefully shutting down for %s...\n", s.String())
		cleanupFunc()
	}()
}
