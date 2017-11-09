package mail

import (
	"bytes"
	"html/template"
	"strconv"

	gomail "gopkg.in/gomail.v2"

	"github.com/JormungandrK/microservice-mail/config"
)

// Info holds info for the email template
type Info struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Email           string `json:"email,omitempty"`
	VerificationURL string `json:"verificationURL,omitempty"`
	Token           string `json:"token,omitempty"`
}

// Send sends an email for verification.
func Send(mailInfo *Info, cfg *config.Config, template string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", cfg.Mail["email"])
	message.SetHeader("To", mailInfo.Email)
	message.SetHeader("Subject", "Verify Your Account!")
	message.SetBody("text/html", template)

	port, err := strconv.Atoi(cfg.Mail["port"])
	if err != nil {
		return err
	}
	d := gomail.NewDialer(cfg.Mail["host"], port, cfg.Mail["user"], cfg.Mail["password"])

	err = d.DialAndSend(message)
	return err
}

// parseTemplate creates a template using emailTemplate.html
func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}

	// Stores the parsed template
	var buff bytes.Buffer

	// Send the parsed template to buff
	err = tmpl.Execute(&buff, data)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}
