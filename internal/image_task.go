package internal

import (
	"github.com/GregoryAlbouy/shrinker/pkg/queue"
	"github.com/streadway/amqp"
)

// NewProducer returns a producer capable of publishing messages onto the image task queue.
func NewProducer() (queue.Producer, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		return queue.Producer{}, err
	}
	defer conn.Close()

	producer, err := queue.NewProducer(conn)
	if err != nil {
		return queue.Producer{}, err
	}
	return producer, nil
}

// Any code required to implement the image processing message queueing would be here.
// For example, we could interface a third party library for and export a custom function.
