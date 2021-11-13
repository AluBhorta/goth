package otputils

import (
	"crypto/rand"
)

const OTP_CHARS = "1234567890"

func GenerateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(OTP_CHARS)
	for i := 0; i < length; i++ {
		buffer[i] = OTP_CHARS[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}
