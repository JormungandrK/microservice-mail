package config

import (
	"encoding/json"
	"io/ioutil"
)

// MailTemplate holds data for mail template
type MailTemplate struct {
	Filename string            `json:"filename"`
	Subject  string            `json:"subject"`
	Data     map[string]string `json:"data"`
}

// Config holds the microservice full configuration.
type Config struct {
	// TemplatesBaseLocation defines base path to mail templates
	TemplateBaseLocation string `json:"templatesBaseLocation"`

	// VerificationURL contains the url to the verification action
	Template map[string]MailTemplate `json:"templates"`

	// Mail is a map of <property>:<value>. For example,
	// "host": "smtp.example.com"
	Mail map[string]string `json:"mail"`

	// RabbitMQ holds information about the rabbitmq server
	AMQPConfig map[string]string `json:"amqpConfig"`
}

// LoadConfig loads a Config from a configuration JSON file.
func LoadConfig(confFile string) (*Config, error) {
	confBytes, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(confBytes, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
