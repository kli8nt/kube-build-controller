package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/Mohamed-Rafraf/rabbitMQ/config"
	"github.com/Mohamed-Rafraf/rabbitMQ/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeployTime() {
	for {
		// Generate a random number between 1 and 100
		randomNumber := rand.Intn(100) + 1

		if randomNumber == 50 {
			log.Println("Publishing To Deploy Queue")
			// Create a struct instance
			deploy := pkg.Deploy{
				AppName: "Hello",
				Port:    42,
			}

			// Serialize the struct to JSON
			jsonData, err := json.Marshal(deploy)
			if err != nil {
				log.Fatalf("Failed to serialize struct to JSON: %v", err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err = config.Channel.PublishWithContext(ctx,
				"",                      // exchange
				config.DeployQueue.Name, // routing key
				false,                   // mandatory
				false,                   // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        jsonData,
				})

			if err != nil {
				log.Println(err)
			}
			cancel()
		}
	}
}
func main() {

	// Consume messages from the queue
	msgs, err := config.Channel.Consume(
		"Build",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
