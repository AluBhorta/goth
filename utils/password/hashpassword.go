package passwordutils

import "golang.org/x/crypto/bcrypt"

func GetHashedPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPass), err
}

func DoesPasswordMatchHash(originalHashedPassword, testPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(originalHashedPassword), []byte(testPassword))
	return err == nil
}
