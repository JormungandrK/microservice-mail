package mail

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strconv"

	gomail "gopkg.in/gomail.v2"

	"github.com/Microkubes/microservice-mail/config"
)

// Info holds info for the email template
type Info struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Email           string `json:"email,omitempty"`
	VerificationURL string `json:"verificationURL,omitempty"`
	Token           string `json:"token,omitempty"`
}

// SMTP Auth to be used with unencrypted connections.
// This MUST NOT be used in production, as it allows sending info and data over
// unencrypted channels.
type unencryptedAuth struct {
	smtp.Auth
}

// Start starts the auth process for the specified SMTP server.
func (u *unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	server.TLS = true
	return u.Auth.Start(server)
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

	// ALLOW_UNENCRYPTED_CONNECTION is intended to be used only when developing and testing
	// with internal/fake SMTP servers when the risk of sending data to the SMTP server
	// over unencrypted channel is negligible. In production this setting must be turned OFF.
	allowUnencryptedConnection := os.Getenv("ALLOW_UNENCRYPTED_CONNECTION")
	if allow, _ := strconv.ParseBool(allowUnencryptedConnection); allow {
		d.Auth = &unencryptedAuth{
			smtp.PlainAuth("", cfg.Mail["user"], cfg.Mail["password"], cfg.Mail["host"]),
		}
		log.Println("[WARN] Authenticating and sending over unencrypted connection is allowed.")
	}

	err = d.DialAndSend(message)
	return err
}

// ParseTemplate creates a template using emailTemplate.html
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
