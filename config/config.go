package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	amqp "github.com/rabbitmq/amqp091-go"
)

type config struct {
	Server   string `env:"HOST" envDefault:"localhost" `
	Port     string `env:"PORT" envDefault:"5672"`
	Password string `env:"PASSWORD" envDefault:"guest"`
	Username string `env:"USERNAME" envDefault:"guest"`
}

var Config config
var Channel *amqp.Channel
var DeployQueue amqp.Queue
var BuildQueue amqp.Queue

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func init() {
	config := config{}
	if err := env.Parse(&config); err != nil {
		fmt.Printf("%+v\n", err)
	}
	log.Println(config)
	url := "amqp://" + config.Username + ":" + config.Password + "@" + config.Server + ":" + config.Port + "/"

	log.Println(url)
	conn, err := amqp.Dial(url)

	failOnError(err, "Failed to connect to RabbitMQ")

	Channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	BuildQueue, err = Channel.QueueDeclare(
		"Build", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a build queue")

	DeployQueue, err = Channel.QueueDeclare(
		"Deploy", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a deploy queue")
	/*
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

		body := "Hello World!"
		err = Channel.PublishWithContext(ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")*/
}
