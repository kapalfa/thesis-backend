package functions

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmail(email string, verificationToken string, emailType string) {
	from := mail.NewEmail("online-code-editor", "androkosm2000@gmail.com")
	to := mail.NewEmail("Example User", email)
	var subject string
	var plainTextContent string
	var htmlContent string

	if emailType == "verify-email" {
		subject = "Verification Email"
		plainTextContent = "Hello and Welcome to our online code editor! You are just one step away from signing up"
		htmlContent = fmt.Sprintf("<h1>Verify your email</h1><p>Use the following code in the confirmation form: %s</p>", verificationToken)
	}
	if emailType == "reset-password" {
		subject = "Reset Password"
		plainTextContent = "Hello, you have requested to reset your password. Use the following code to reset your password"
		htmlContent = fmt.Sprintf("<h1>Reset your password</h1><p>Use the following code to reset your password: %s</p>", verificationToken)
	}

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
