package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Mohamed-Rafraf/kube-builder/config"
	"github.com/Mohamed-Rafraf/kube-builder/pkg"
	"github.com/Mohamed-Rafraf/kube-builder/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

var err error

func DeployTime() {
	for {
		var status *pkg.Status
		// Generate a random number between 1 and 100
		status, err = utils.Delete()
		if err != nil {
			log.Println(err)
			continue
		}
		if status != nil && status.Status == "Succeeded" {
			log.Println("Publishing To Deploy Queue")

			// Serialize the struct to JSON
			jsonData, err := json.Marshal(status)
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
	var data pkg.Data
	go func() {
		for d := range msgs {
			err = json.Unmarshal(d.Body, &data)
			log.Println("Data Recieved", data)
			if err != nil {
				log.Println(err)
				continue
			}

			err = utils.Build(&data)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
