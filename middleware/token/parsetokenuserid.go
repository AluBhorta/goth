package tokenmiddleware

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	commonmodels "github.com/alubhorta/goth/models/common"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func ParseTokenUserId(c *fiber.Ctx) error {
	authHeader := c.Request().Header.Peek("Authorization")
	authHeaderCopy := make([]byte, len(authHeader))
	copy(authHeaderCopy, authHeader)
	authHeaderStr := string(authHeaderCopy)

	splitted := strings.Split(authHeaderStr, " ")
	if len(splitted) != 2 {
		msg := "invalid token provided."
		log.Println(msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	accessToken := splitted[1]
	// NOTE: possible to refactor token parsing from header into tokenutils func
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		signingKey := os.Getenv("ACCESS_TOKEN_SIGNING_KEY")
		return []byte(signingKey), nil
	})
	if err != nil {
		msg := "failed to parse or validate token."
		log.Println(msg, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		msg := "invalid token or claim typecast error."
		log.Println(msg, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	userId, ok := claims["userId"].(string)
	if !ok || len(userId) <= 0 {
		msg := "invalid user id provided in claim."
		log.Println(msg, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	prevCtx := c.UserContext().Value(commonmodels.CommonCtx{}).(*commonmodels.CommonCtx)
	newCtx := context.WithValue(
		context.Background(),
		commonmodels.CommonCtx{},
		&commonmodels.CommonCtx{
			Clients: prevCtx.Clients,
			UserId:  userId,
		},
	)
	c.SetUserContext(newCtx)

	return c.Next()
}
