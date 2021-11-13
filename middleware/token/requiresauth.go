package tokenmiddleware

import (
	"log"
	"os"
	"strings"

	customerrors "github.com/alubhorta/goth/custom/errors"
	commonmodels "github.com/alubhorta/goth/models/common"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func RequiresAuth(c *fiber.Ctx) error {
	authHeader := c.Request().Header.Peek("Authorization")
	authHeaderCopy := make([]byte, len(authHeader))
	copy(authHeaderCopy, authHeader)
	authHeaderStr := string(authHeaderCopy)

	splitted := strings.Split(authHeaderStr, " ")
	if len(splitted) == 2 {
		accessToken := splitted[1]

		cc := c.UserContext().Value(commonmodels.CommonCtx{}).(*commonmodels.CommonCtx).Clients
		cacheClient := cc.CacheClient

		res, err := cacheClient.Get(accessToken)
		if err != nil && err != customerrors.ErrNotFound {
			msg := "failed to lookup cache."
			log.Println(msg, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
		} else if res == "blacklist:access" {
			msg := "blacklisted token used."
			log.Println(msg)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": msg, "payload": nil})
		}
	} else {
		msg := "invalid token provided."
		log.Println(msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	return jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("ACCESS_TOKEN_SIGNING_KEY")),
	})(c)
}
