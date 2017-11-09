package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/JormungandrK/microservice-mail/config"
	"github.com/JormungandrK/microservice-mail/mail"
	"github.com/JormungandrK/microservice-tools/rabbitmq"
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
	cf := os.Getenv("SERVICE_CONFIG_FILE")
	if cf == "" {
		cf = "/run/secrets/microservice_mail_config.json"
	}
	cfg, err := config.LoadConfig(cf)
	logOnError(err, "Failed to read the config file!")

	conn, ch, err := rabbitmq.Dial(
		cfg.RabbitMQ["username"],
		cfg.RabbitMQ["password"],
		cfg.RabbitMQ["host"],
		cfg.RabbitMQ["post"],
	)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	defer ch.Close()

	channel := &rabbitmq.AMQPChannel{ch}
	msgs, err := channel.Receive("verification-email")
	if err != nil {
		logOnError(err, "Failed to consume the channel")
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			mailInfo := mail.Info{}
			err := json.Unmarshal(d.Body, &mailInfo)
			logOnError(err, "Failed to unmarshal body")

			mailInfo.VerificationURL = cfg.VerificationURL
			template, err := mail.ParseTemplate("./public/mail/template.html", mailInfo)
			if err != nil {
				logOnError(err, "Failed to parse mail tamplate")
			}

			err = mail.Send(&mailInfo, cfg, template)
			logOnError(err, fmt.Sprintf("Failed to send mail to %s", mailInfo.Email))

			d.Ack(false)
			log.Printf("Done")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	// Make goroutine to work forever
	<-forever
}
