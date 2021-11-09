package tokenutils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateNewAccessToken(userId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	maxAgeInSeconds, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_MAX_AGE_IN_SECONDS"))
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = userId
	// claims["jti"] = fmt.Sprintf("%v", uuid.New())
	now := time.Now()
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(time.Second * time.Duration(maxAgeInSeconds)).Unix()

	signingKey := os.Getenv("ACCESS_TOKEN_SIGNING_KEY")
	return token.SignedString([]byte(signingKey))
}

func CreateNewRefreshToken(userId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	maxAgeInSeconds, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_MAX_AGE_IN_SECONDS"))
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = userId
	// claims["jti"] = fmt.Sprintf("%v", uuid.New())
	now := time.Now()
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(time.Second * time.Duration(maxAgeInSeconds)).Unix()

	signingKey := os.Getenv("REFRESH_TOKEN_SIGNING_KEY")
	return token.SignedString([]byte(signingKey))
}
