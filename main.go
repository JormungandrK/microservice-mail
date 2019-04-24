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

func logOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {

	cfg := getConfig()

	rabbitMQChannel := getRabbitMQChannel(cfg)
	channel := rabbitmq.AMQPChannel{
		Channel: rabbitMQChannel,
	}

	for {
		deliveryList, err := channel.Receive("email")
		logOnError(err, "Failed to consume the channel")

		for delivery := range deliveryList {
			go handleDelivery(delivery)
		}
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
}

func handleDelivery(delivery Delivery) {
	log.Printf("Received a message: %s", delivery.Body)
	message, err := mail.ParseRabbitMQMessage(&delivery.Body)
	logOnError(err, err.Error())
	body, err := mail.GenerateMailBody(cfg, &message)
	logOnError(err, err.Error())
	err = mail.SendMail(&message, cfg, &body)
	logOnError(err, fmt.Sprintf("Failed to send mail to %s", message.Email))
	delivery.Ack(false)
	log.Printf("Done")
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

func getRabbitMQChannel(cfg *config.Config) *amqp.Channel {
	conn, ch, err := rabbitmq.Dial(
		cfg.RabbitMQ["username"],
		cfg.RabbitMQ["password"],
		cfg.RabbitMQ["host"],
		cfg.RabbitMQ["post"],
	)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	defer ch.Close()
	return ch
}
