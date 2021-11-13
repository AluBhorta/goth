package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	authapi "github.com/alubhorta/goth/api/auth"
	userapi "github.com/alubhorta/goth/api/user"
	"github.com/alubhorta/goth/db/cacheclient"
	"github.com/alubhorta/goth/db/dbclient"
	tokenmw "github.com/alubhorta/goth/middleware/token"
	commonmodels "github.com/alubhorta/goth/models/common"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app := fiber.New()

	app.Use(cors.New())

	dbclient := &dbclient.MongoDbClient{}
	dbclient.Init()

	redisClient := &cacheclient.RedisClient{}
	redisClient.Init()

	commonClients := &commonmodels.CommonClients{
		DbClient:    dbclient,
		CacheClient: redisClient,
	}
	userCtx := context.WithValue(
		context.Background(),
		commonmodels.CommonCtx{},
		&commonmodels.CommonCtx{Clients: commonClients, UserId: ""},
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
	app.Delete("/api/v1/auth/delete", tokenmw.ParseTokenUserId, tokenmw.RequiresAuth, authapi.DeleteAccount)

	// user routes
	app.Get("/api/v1/user", tokenmw.ParseTokenUserId, tokenmw.RequiresAuth, userapi.GetOne)
	app.Put("/api/v1/user", tokenmw.ParseTokenUserId, tokenmw.RequiresAuth, userapi.UpdateOne)
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
