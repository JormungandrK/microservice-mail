package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Microkubes/microservice-mail/config"
	"github.com/Microkubes/microservice-mail/mail"
	"github.com/Microkubes/microservice-tools/rabbitmq"
	"github.com/streadway/amqp"
)

func logOnError(err error, msg string) bool {
	if err != nil {
		fmt.Println(msg, err)
		return true
	}
	return false
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf(msg, err)
		panic(fmt.Sprintf(msg, err))
	}
}

func main() {

	cfg := getConfig()

	amqpChannel := getAMQPChannel(cfg)
	channel := rabbitmq.AMQPChannel{
		Channel: amqpChannel,
	}

	for {
		deliveryList, err := channel.Receive("email")
		logOnError(err, "Failed to consume the channel")
		for delivery := range deliveryList {
			go handleDelivery(delivery, cfg)
		}
	}
}

func handleDelivery(delivery amqp.Delivery, cfg *config.Config) bool {

	log.Printf("Received a message: %s", delivery.Body)

	message, err := mail.ParseAMQPMessage(&delivery.Body)
	if logOnError(err, "Failed to parse AMQP Message") {
		delivery.Ack(false)
		return false
	}

	body, err := mail.GenerateMailBody(cfg, &message)
	if logOnError(err, "Failed to generate mail body for template "+message.TemplateName) {
		delivery.Ack(false)
		return false
	}

	err = mail.SendMail(&message, cfg, &body)
	if logOnError(err, "Failed to send mail to "+message.Email) {
		delivery.Ack(false)
		return false
	}

	log.Printf("Message to " + message.Email + " sucessfully sended!")
	delivery.Ack(false)
	return true
}

func getConfig() *config.Config {
	cf := os.Getenv("SERVICE_CONFIG_FILE")
	if cf == "" {
		cf = "/run/secrets/microservice_mail_config.json"
	}
	cfg, err := config.LoadConfig(cf)
	logOnError(err, "Failed to read the config file!")
	return cfg
}

func getAMQPChannel(cfg *config.Config) *amqp.Channel {
	_, ch, err := rabbitmq.Dial(
		cfg.RabbitMQ["username"],
		cfg.RabbitMQ["password"],
		cfg.RabbitMQ["host"],
		cfg.RabbitMQ["port"],
	)
	failOnError(err, "Failed to connect to RabbitMQ")
	return ch
}
