package authapi

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	customerrors "github.com/alubhorta/goth/custom/errors"
	authmodels "github.com/alubhorta/goth/models/auth"
	commonclients "github.com/alubhorta/goth/models/common"
	usermodels "github.com/alubhorta/goth/models/user"
	otputils "github.com/alubhorta/goth/utils/otp"
	passwordutils "github.com/alubhorta/goth/utils/password"
	tokenutils "github.com/alubhorta/goth/utils/token"
	validationutils "github.com/alubhorta/goth/utils/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
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

	// generate new token pair
	accessToken, err := tokenutils.CreateNewAccessToken(userId)
	if err != nil {
		msg := "failed to generate access token."
		log.Println(msg, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	}
	refreshToken, err := tokenutils.CreateNewRefreshToken(userId)
	if err != nil {
		msg := "failed to generate refresh token."
		log.Println(msg, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	msg := "successful signup completed."
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

func Login(c *fiber.Ctx) error {
	input := new(authmodels.LoginInput)
	if err := c.BodyParser(input); err != nil {
		msg := "invalid input."
		log.Println(msg, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if !validationutils.IsValidEmail(input.Email) {
		msg := "invalid email provided."
		log.Println(msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	cc := c.UserContext().Value(commonclients.CommonClients{}).(*commonclients.CommonClients)
	dbclient := cc.DbClient

	authCred, err := dbclient.AuthAccess.GetAuthCredentialByEmail(input.Email)
	if err == customerrors.ErrNotFound || (err == nil && authCred == nil) {
		msg := "no such user found."
		log.Println(msg)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if err != nil {
		msg := "failed to login."
		log.Println(msg, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	}
	matches := passwordutils.DoesPasswordMatchHash(authCred.HashedPassword, input.Password)
	if !matches {
		msg := "invalid password provided."
		log.Println(msg, "input password does not match hashed password")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	accessToken, err := tokenutils.CreateNewAccessToken(authCred.UserId)
	if err != nil {
		msg := "failed to generate access token."
		log.Println(msg, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	}
	refreshToken, err := tokenutils.CreateNewRefreshToken(authCred.UserId)
	if err != nil {
		msg := "failed to generate refresh token."
		log.Println(msg, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	msg := "successfully logged in user."
	log.Println(msg, authCred.UserId)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": msg,
		"payload": fiber.Map{
			"userId": authCred.UserId,
			"tokens": fiber.Map{
				"access":  accessToken,
				"refresh": refreshToken,
			},
		},
	})
}

func Logout(c *fiber.Ctx) error {
	input := new(authmodels.LogoutInput)
	if err := c.BodyParser(input); err != nil {
		msg := "invalid input."
		log.Println(msg, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	if input.AccessToken == "" && input.RefreshToken == "" {
		msg := "invalid tokens provided."
		log.Println(msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	} else {
		accessMaxAgeInSeconds, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_MAX_AGE_IN_SECONDS"))
		if err != nil {
			msg := "error in type conversion."
			log.Println(msg, "from string to int", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
		}
		refreshMaxAgeInSeconds, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_MAX_AGE_IN_SECONDS"))
		if err != nil {
			msg := "error in type conversion."
			log.Println(msg, "from string to int", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
		}

		cc := c.UserContext().Value(commonclients.CommonClients{}).(*commonclients.CommonClients)
		cacheClient := cc.CacheClient

		cacheClient.Set(input.AccessToken, "blacklist:access", time.Second*time.Duration(accessMaxAgeInSeconds))
		cacheClient.Set(input.RefreshToken, "blacklist:refresh", time.Second*time.Duration(refreshMaxAgeInSeconds))

		msg := "successfully logged out."
		log.Println(msg)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": msg, "payload": nil})
	}
}

func Refresh(c *fiber.Ctx) error {
	input := new(authmodels.RefreshInput)
	if err := c.BodyParser(input); err != nil {
		msg := "invalid input."
		log.Println(msg, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	cc := c.UserContext().Value(commonclients.CommonClients{}).(*commonclients.CommonClients)
	cacheClient := cc.CacheClient

	res, err := cacheClient.Get(input.RefreshToken)
	if err != nil && err != customerrors.ErrNotFound {
		msg := "failed to lookup cache."
		log.Println(msg, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if res == "blacklist:refresh" {
		msg := "blacklisted token used."
		log.Println(msg)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		signingKey := os.Getenv("REFRESH_TOKEN_SIGNING_KEY")
		return []byte(signingKey), nil
	})
	if err != nil {
		msg := "failed to parse or validate token."
		log.Println(msg, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId, ok := claims["userId"].(string)
		if !ok || len(userId) <= 0 {
			msg := "invalid user id provided in claim."
			log.Println(msg, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
		}

		accessToken, err := tokenutils.CreateNewAccessToken(userId)
		if err != nil {
			msg := "failed to generate access token."
			log.Println(msg, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
		}
		refreshToken, err := tokenutils.CreateNewRefreshToken(userId)
		if err != nil {
			msg := "failed to generate refresh token."
			log.Println(msg, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
		}

		msg := "successfully refreshed tokens."
		log.Println(msg, "for userId: ", userId)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": msg,
			"payload": fiber.Map{
				"tokens": fiber.Map{
					"access":  accessToken,
					"refresh": refreshToken,
				},
			},
		})
	} else {
		msg := "invalid token or claim typecast error."
		log.Println(msg, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	}
}

func ResetPasswordInit(c *fiber.Ctx) error {
	// parse email input from request
	input := new(authmodels.ResetInitInput)
	err := c.BodyParser(input)
	if err != nil || input.Email == "" {
		msg := "invalid input."
		log.Println(msg, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if !validationutils.IsValidEmail(input.Email) {
		msg := "invalid email provided."
		log.Println(msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	cc := c.UserContext().Value(commonclients.CommonClients{}).(*commonclients.CommonClients)
	dbClient := cc.DbClient

	// send 404 if email doesn't exist
	_, err = dbClient.AuthAccess.GetAuthCredentialByEmail(input.Email)
	if err == customerrors.ErrNotFound {
		msg := "email does not exist."
		log.Println(msg, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if err != nil {
		msg := "failed to read from database."
		log.Println(msg, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	cacheClient := cc.CacheClient
	cacheKey := "resetOTP:" + input.Email
	exists, err := cacheClient.Exists(cacheKey)
	if err != nil {
		msg := "failed to read cache."
		log.Println(msg, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if exists {
		msg := "password reset already initiated for this email. check your email or try after 2 minutes."
		log.Println(msg)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": msg, "payload": nil})
	} // else carry on

	// save otp to cache with TTL=2min
	otp, err := otputils.GenerateOTP(6)
	if err != nil {
		msg := "failed to generate otp."
		log.Println(msg, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	}
	cacheClient.Set(cacheKey, otp, time.Second*120)
	log.Println("otp generated:", otp, "\tcacheKey: ", cacheKey)

	// TODO: send otp via email to input.Email

	msg := "a verification code (otp) is sent to your email. reset your password within the next 2 minutes."
	log.Println(msg)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": msg,
		"payload": nil,
	})
}

func ResetPasswordVerify(c *fiber.Ctx) error {
	input := new(authmodels.ResetVerifyInput)
	err := c.BodyParser(input)
	if err != nil || input.Email == "" || input.Otp == "" || input.NewPassword == "" {
		msg := "invalid input."
		log.Println(msg, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if !validationutils.IsValidEmail(input.Email) {
		msg := "invalid email provided."
		log.Println(msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if len(input.NewPassword) < 6 {
		msg := "invalid input - password must be at least 6 characters."
		log.Println(msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	cacheKey := "resetOTP:" + input.Email

	cc := c.UserContext().Value(commonclients.CommonClients{}).(*commonclients.CommonClients)
	cacheClient := cc.CacheClient

	val, err := cacheClient.Get(cacheKey)
	if err == customerrors.ErrNotFound {
		msg := "not found - invalid input or expired key."
		log.Println(msg)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if err != nil {
		msg := "failed to read from cache."
		log.Println(msg, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if val != input.Otp {
		msg := "invalid input - otp mismatch."
		log.Println(msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	} // else all good, change the password

	newHasedPass, err := passwordutils.GetHashedPassword(input.NewPassword)
	if err != nil {
		msg := "could not hash password."
		log.Println(msg, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	dbclient := cc.DbClient
	err = dbclient.AuthAccess.UpdateUserAuthPassword(input.Email, newHasedPass)
	if err == customerrors.ErrNotFound {
		msg := "no such user found."
		log.Println(msg, "with email:", input.Email)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if err != nil {
		msg := "failed to update password."
		log.Println(msg, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	msg := "password successfully reset."
	log.Println(msg, "for user with email:", input.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": msg,
		"payload": nil,
	})
}

func DeleteAccount(c *fiber.Ctx) error {
	userId := c.Params("id")
	if userId == "" {
		msg := "empty user id provided."
		log.Println(msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	cc := c.UserContext().Value(commonclients.CommonClients{}).(*commonclients.CommonClients)
	dbclient := cc.DbClient

	// TODO: [transaction safety] - find out a way to delete  both documents atomically
	err := dbclient.AuthAccess.DeleteAnAuthCredential(userId)
	if err == customerrors.ErrNotFound {
		msg := "no such user credential found for deletion."
		log.Println(msg, err, "id:", userId)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if err != nil {
		msg := "failed to delete user credential."
		log.Println(msg, err, "id:", userId)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	err = dbclient.UserAccess.DeleteAUser(userId)
	if err == customerrors.ErrNotFound {
		msg := "no such user found for deletion."
		log.Println(msg, err, "id:", userId)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": msg, "payload": nil})
	} else if err != nil {
		msg := "failed to delete user."
		log.Println(msg, "id:", userId, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": msg, "payload": nil})
	}

	msg := "successfully deleted user."
	log.Println(msg, "id:", userId)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": msg, "payload": nil})
}
