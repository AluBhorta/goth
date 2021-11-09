package authapi

import (
	"fmt"
	"log"
	"time"

	authmodels "github.com/alubhorta/goth/models/auth"
	usermodels "github.com/alubhorta/goth/models/user"
	passwordutils "github.com/alubhorta/goth/utils/password"
	validationutils "github.com/alubhorta/goth/utils/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Signup(c *fiber.Ctx) error {
	// / input validation
	input := new(authmodels.SignupInput)
	if err := c.BodyParser(input); err != nil {
		msg := "invalid input"
		log.Println(msg, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if input.Email == "" || input.Password == "" || input.FirstName == "" || input.LastName == "" {
		msg := "invalid input - missing required fields"
		log.Println(msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if !validationutils.IsValidEmail(input.Email) {
		msg := "invalid email provided"
		log.Println(msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if len(input.Password) < 6 {
		msg := "invalid input - password must be at least 6 characters"
		log.Println(msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	// hash password
	hasedPass, err := passwordutils.GetHashedPassword(input.Password)
	if err != nil {
		msg := "could not hash password"
		log.Println(msg, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	userId := fmt.Sprintf("%v", uuid.New())
	now := time.Now()
	authCred := &authmodels.UserAuthCredential{
		Email:          input.Email,
		UserId:         userId,
		HashedPassword: hasedPass,
		CreatedAt:      now,
		ModifiedAt:     now,
	}

	// TODO: save new auth credential
	log.Println("authCred", authCred)

	createUserInput := &usermodels.CreateUserInfoInput{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}
	// TODO: create and save new User
	log.Println("createUserInput", createUserInput)

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
