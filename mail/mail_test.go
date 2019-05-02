package mail

import (
	"encoding/json"
	"testing"

	"github.com/Microkubes/microservice-mail/config"
	"github.com/Microkubes/microservice-tools/rabbitmq"
	"github.com/streadway/amqp"
)

func getConfig() *config.Config {
	confBytes := []byte(`{
		"templatesBaseLocation": "./",
		"templates": {
			"sometemplate": {
				"filename" : "test-template.html",
				"subject": "Template header",
				"data": {
					"test": "value"
				}
			}
		},
		"mail": {
			"host": "smtp.mailtrap.io",
			"port": "2525",
			"user": "7dfa3710bee1c3",
			"password": "68d8ccb96fb52b",
			"email": "e0e3decc9e-f10431@inbox.mailtrap.io"
		}
	}`)
	config := &config.Config{}
	json.Unmarshal(confBytes, config)
	return config
}

func getAMQPChannel(cfg *config.Config) *amqp.Channel {
	_, ch, _ := rabbitmq.Dial(
		cfg.AMQPConfig["username"],
		cfg.AMQPConfig["password"],
		cfg.AMQPConfig["host"],
		cfg.AMQPConfig["port"],
	)
	return ch
}

func getAMQPMessage() []byte {
	messageBody := []byte(`{
		"email" : "kalevski@keitaro.com",
		"template": "sometemplate",
		"data": { 
			"test": "value"
			}
		}`)
	return messageBody
}

func TestParseAMQPMessage(t *testing.T) {
	messageBody := getAMQPMessage()
	parsed, err := ParseAMQPMessage(&messageBody)
	if err != nil {
		t.Errorf("Can't parse AMQP Message")
	}
	if parsed.Email != "kalevski@keitaro.com" {
		t.Errorf("Something went wrong while parsing email property")
	}
	if parsed.TemplateName != "sometemplate" {
		t.Errorf("Something went wrong while parsing template property")
	}
	if parsed.Data["test"] != "value" {
		t.Errorf("Something went wrong while parsing data map")
	}
}

func TestGenerateMailBody(t *testing.T) {
	cfg := getConfig()
	message := getAMQPMessage()
	amqpMessage, _ := ParseAMQPMessage(&message)
	body, err := GenerateMailBody(cfg, &amqpMessage)
	if err != nil {
		t.Errorf("Failed while generating mail content" + err.Error())
	}
	if body == "" {
		t.Errorf("Wrong generated content")
	}
}

func TestSendMail(t *testing.T) {
	cfg := getConfig()
	message := getAMQPMessage()
	amqpMessage, _ := ParseAMQPMessage(&message)
	body, _ := GenerateMailBody(cfg, &amqpMessage)
	SendMail(&amqpMessage, cfg, &body)
}
