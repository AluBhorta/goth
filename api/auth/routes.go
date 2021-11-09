package authapi

import (
	"fmt"
	"log"
	"time"

	customerrors "github.com/alubhorta/goth/custom/errors"
	authmodels "github.com/alubhorta/goth/models/auth"
	commonclients "github.com/alubhorta/goth/models/common"
	usermodels "github.com/alubhorta/goth/models/user"
	passwordutils "github.com/alubhorta/goth/utils/password"
	validationutils "github.com/alubhorta/goth/utils/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Signup(c *fiber.Ctx) error {
	// input validation
	input := new(authmodels.SignupInput)
	if err := c.BodyParser(input); err != nil {
		msg := "invalid input."
		log.Println(msg, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if input.Email == "" || input.Password == "" || input.FirstName == "" || input.LastName == "" {
		msg := "invalid input - missing required fields."
		log.Println(msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if !validationutils.IsValidEmail(input.Email) {
		msg := "invalid email provided."
		log.Println(msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if len(input.Password) < 6 {
		msg := "invalid input - password must be at least 6 characters."
		log.Println(msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	// hash password
	hasedPass, err := passwordutils.GetHashedPassword(input.Password)
	if err != nil {
		msg := "could not hash password."
		log.Println(msg, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	cc := c.UserContext().Value(commonclients.CommonClients{}).(*commonclients.CommonClients)
	dbclient := cc.DbClient

	// create auth credential
	userId := fmt.Sprintf("%v", uuid.New())
	now := time.Now()
	authCred := &authmodels.UserAuthCredential{
		Email:          input.Email,
		UserId:         userId,
		HashedPassword: hasedPass,
		CreatedAt:      now,
		ModifiedAt:     now,
	}
	err = dbclient.AuthAccess.CreateNewUserAuthCredential(authCred)
	if err == customerrors.ErrDuplicateKey {
		msg := "failed to create auth credentials - duplicate key."
		log.Println(msg, err)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if err != nil {
		msg := "failed to create auth credentials."
		log.Println(msg, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	// create minimal model for userInfo
	createUserInput := &usermodels.CreateUserInfoInput{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}
	err = dbclient.UserAccess.CreateAUser(userId, createUserInput)
	if err == customerrors.ErrDuplicateKey {
		msg := "failed to create user - duplicate key."
		log.Println(msg, err)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if err != nil {
		msg := "failed to create user."
		log.Println(msg, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	// TODO: generate new token pair
	accessToken := "mock-accessToken"
	refreshToken := "mock-refreshToken"

	msg := "successful signup up completed."

	log.Println(msg, "userId:", userId)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": msg,
		"payload": fiber.Map{
			"userId": userId,
			"tokens": fiber.Map{
				"access":  accessToken,
				"refresh": refreshToken,
			},
		},
	})
}

func Login(c *fiber.Ctx) error { return nil }

func Logout(c *fiber.Ctx) error { return nil }

func Refresh(c *fiber.Ctx) error { return nil }

func ResetPasswordInit(c *fiber.Ctx) error { return nil }

func ResetPasswordVerify(c *fiber.Ctx) error { return nil }

func DeleteAccount(c *fiber.Ctx) error { return nil }
