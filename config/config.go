package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
	amqp "github.com/rabbitmq/amqp091-go"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type configuration struct {
	Server   string `env:"HOST" envDefault:"localhost" `
	Port     string `env:"PORT" envDefault:"5672"`
	Password string `env:"PASSWORD" envDefault:"guest"`
	Username string `env:"USERNAME" envDefault:"guest"`
}

var Config configuration
var Channel *amqp.Channel
var DeployQueue amqp.Queue
var BuildQueue amqp.Queue
var StatusQueue amqp.Queue
var err error
var Clientset *kubernetes.Clientset
var config *rest.Config

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func init() {
	err = InitConfig()
	if err != nil {
		log.Println(err)
	}

	err = ClusterConfig()
	if err != nil {
		log.Println(err)
	}

	InitRabbit()

}

func ClusterConfig() error {
	var err error

	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")

	flag.Parse()

	if *kubeconfig != "" {
		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			return err
		}

		// create the Clientset
		Clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			return err
		}

	} else {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			return err
		}
		// creates the Clientset
		Clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			return err
		}

	}

	return nil
}

func InitConfig() error {
	err := env.Parse(&Config)
	if err != nil {
		return err
	}
	return nil

}

func InitRabbit() {
	url := "amqp://" + Config.Username + ":" + Config.Password + "@" + Config.Server + ":" + Config.Port + "/"

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

	StatusQueue, err = Channel.QueueDeclare(
		"Status", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a status queue")
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
