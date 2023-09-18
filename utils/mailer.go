package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	smtpHost string
	smtpPort string
	from     string
	to       string
	username string
	password string
)

const (
	UPDATES_TEMPLATE = "-- '%s' by %s\n"
)

func init() {
	err := godotenv.Load("mail.env")
	if err != nil {
		log.Panicln("Something went wrong while loading environment variables")
	}

	smtpHost = strings.TrimSpace(os.Getenv("SMTP_HOST"))
	smtpPort = strings.TrimSpace(os.Getenv("SMTP_PORT"))
	from = strings.TrimSpace(os.Getenv("MAIL_FROM"))
	to = strings.TrimSpace(os.Getenv("MAIL_TO"))
	password = os.Getenv("MAIL_PASSWORD")
	// fmt.Printf("Host: %s\nPort: %s\nFrom: %s\nTo: %s\nPassword: %s\n\n", smtpHost, smtpPort, from, to, password)
}

func SendUpdateMail(updates map[string]string) {
	fmt.Printf("Host: %s\nPort: %s\nFrom: %s\nTo: %s\nPassword: %s\n\n", smtpHost, smtpPort, from, to, password)
	message := []byte(generateMailBody(updates))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		log.Panicf("Something went wrong while sending mail: %s\n", err.Error())
	}
}

func generateMailBody(updates map[string]string) string {
	updateString := ""

	for author, title := range updates {
		updateString += fmt.Sprintf(UPDATES_TEMPLATE, title, author)
	}

	mailBody := `Hi,

    The following feeds have been updated:
    %s

    Regards`

	return fmt.Sprintf(mailBody, updateString)
}
