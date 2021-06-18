package main

import (
	"log"

	"github.com/GregoryAlbouy/shrinker/pkg/dotenv"
	"github.com/GregoryAlbouy/shrinker/pkg/queue"
	"github.com/streadway/amqp"
)

const defaultEnvPath = "./.env"

var env = map[string]string{
	"QUEUE_URL": "",
}

func main() {
	envPath := dotenv.GetPath(defaultEnvPath)

	if err := dotenv.Load(envPath, &env); err != nil {
		log.Fatal(err)
	}

	conn, err := amqp.Dial(env["QUEUE_URL"])
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	consumer, err := queue.NewConsumer(conn)
	if err != nil {
		log.Fatal(err)
	}
	consumer.Listen(logMessage)
}

func logMessage(d amqp.Delivery) error {
	log.Printf("received a message: %s", d.Body)
	return nil
}
