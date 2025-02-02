package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strconv"

	"github.com/Microkubes/microservice-mail/config"
	"gopkg.in/gomail.v2"
)

// AMQPMessage message received from AMQP server
type AMQPMessage struct {
	Email        string            `json:"email,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
	TemplateName string            `json:"template,omitempty"`
}

// VerificationMail object specified for verification mails
type VerificationMail struct {
	URL   string `json:"url,omitempty"`
	Token string `json:"token,omitempty"`
}

// ForgotPasswordMail object specified for forgot password mails
type ForgotPasswordMail struct {
	URL  string `json:"url,omitempty"`
	Code string `json:"code,omitempty"`
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

// ParseAMQPMessage helper for parsing AMQP message Body to our AMQPMessage
func ParseAMQPMessage(body *[]byte) (AMQPMessage, error) {
	msg := AMQPMessage{}
	err := json.Unmarshal(*body, &msg)
	if err != nil {
		return msg, fmt.Errorf("Failed to parse message from AMQP Server")
	}
	if msg.Email == "" {
		return msg, fmt.Errorf("Every AMQP Message must contain Email property")
	}
	if msg.TemplateName == "" {
		return msg, fmt.Errorf("Every AMQP Message must contain TemplateName property")
	}
	return msg, nil
}

// GenerateMailBody generates mail body from configuration & message
func GenerateMailBody(cfg *config.Config, message *AMQPMessage) (string, error) {
	templateConfig, success := cfg.Template[message.TemplateName]

	if !success {
		return "", fmt.Errorf("Doesn't exist template config for received message [" + message.TemplateName + "]")
	}
	if message.Data == nil {
		message.Data = map[string]string{}
	}
	for key, value := range templateConfig.Data {
		message.Data[key] = value
	}

	content, err := parseTemplate(cfg.TemplateBaseLocation, templateConfig.Filename, message.Data)
	if err != nil {
		return "", fmt.Errorf("Failed on parsing template: %s", err)
	}
	return content, nil
}

// SendMail sends an email for verification.
func SendMail(message *AMQPMessage, cfg *config.Config, body *string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", cfg.Mail["email"])
	msg.SetHeader("To", message.Email)
	msg.SetHeader("Subject", cfg.Template[message.TemplateName].Subject)
	msg.SetBody("text/html", *body)

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

	err = d.DialAndSend(msg)
	return err
}

func parseTemplate(baseLocation string, templateFilename string, data interface{}) (string, error) {
	template, err := template.ParseFiles(baseLocation + templateFilename)
	if err != nil {
		return "", err
	}
	var buff bytes.Buffer
	err = template.Execute(&buff, data)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}
