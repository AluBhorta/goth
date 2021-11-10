package userapi

import (
	"log"

	customerrors "github.com/alubhorta/goth/custom/errors"
	commonclients "github.com/alubhorta/goth/models/common"
	usermodels "github.com/alubhorta/goth/models/user"

	"github.com/gofiber/fiber/v2"
)

func GetOne(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		msg := "invalid user Id provided."
		log.Println(msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	cc := c.UserContext().Value(commonclients.CommonClients{}).(*commonclients.CommonClients)
	dbclient := cc.DbClient

	aUser, err := dbclient.UserAccess.GetAUser(id)
	if err == customerrors.ErrNotFound {
		msg := "user does not exist."
		log.Println(msg, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if err != nil {
		msg := "failed to get user."
		log.Println(msg, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	msg := "successfully retrieved user."
	log.Println(msg)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": msg, "payload": fiber.Map{"userInfo": aUser}})
}

func UpdateOne(c *fiber.Ctx) error {
	// NOTE: email is not updateable via this API

	userId := c.Params("id")
	if userId == "" {
		msg := "empty user Id provided."
		log.Println(msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	input := new(usermodels.UpdateUserInfoInput)
	if err := c.BodyParser(input); err != nil {
		msg := "invalid input provided."
		log.Println(msg, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if input.FirstName == "" || input.LastName == "" {
		msg := "invalid input - required fields cannot be empty."
		log.Println(msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	log.Printf("update user with id=%v and updateUserInput=%v\n", userId, input)

	cc := c.UserContext().Value(commonclients.CommonClients{}).(*commonclients.CommonClients)
	dbclient := cc.DbClient

	err := dbclient.UserAccess.UpdateAUser(userId, input)
	if err == customerrors.ErrDuplicateKey {
		msg := "failed to update user - duplicate key."
		log.Println(msg, err)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if err == customerrors.ErrNotFound {
		msg := "no such user found."
		log.Println(msg, err, "id: ", userId)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if err != nil {
		msg := "failed to update user."
		log.Println(msg, err, "id: ", userId)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	msg := "successfully updated user"
	log.Println(msg, "id: ", userId)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": msg, "payload": nil})
}
