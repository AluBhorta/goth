package emailutils

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendMail(toEmail string, fromEmail, subject, body string) error {
	fromName := strings.Split(fromEmail, "@")[0]
	from := mail.NewEmail(fromName, fromEmail)

	toName := strings.Split(toEmail, "@")[0]
	to := mail.NewEmail(toName, toEmail)

	message := mail.NewSingleEmail(from, subject, to, "", body)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)

	if err != nil {
		return err
	} else if response.StatusCode != 202 {
		log.Println("StatusCode: ", response.StatusCode)
		log.Println("Body: ", response.Body)
		log.Println("Headers: ", response.Headers)
		return errors.New(response.Body)
	} else {
		return nil
	}
}

func SendOtpMail(toEmail, otp string) error {
	// NOTE: feel free to update the email composition as per your requirements
	subject := "Password reset code | GOTH"
	fromEmail := os.Getenv("FROM_EMAIL_ADDRESS")
	htmlBody := "<p>Your password reset code is: <strong>" + otp + "</strong></p>"

	return SendMail(toEmail, fromEmail, subject, htmlBody)
}
