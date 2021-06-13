package main

import (
	"log"

	"github.com/GregoryAlbouy/shrinker/pkg/queue"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatalf("worker error: %s", err)
	}
	defer conn.Close()

	consumer, err := queue.NewConsumer(conn)
	if err != nil {
		log.Fatalf("worker error: %s", err)
	}
	consumer.Listen(logMessage)
}

func logMessage(d amqp.Delivery) error {
	log.Printf("Received a message: %s", d.Body)
	return nil
}
