package queue

import "github.com/streadway/amqp"

var queueName string

// SetQueueName sets the queue name used inside the package `queue`.
//
// If you want to use a named queue, you must invoke SetQueueName with the
// name you desire before ever calling Producer.Publish or Consumer.Listen.
//
// If SetQueueName is never invoked, the package use the empty string ""
// as the queue name. A nameless queue is assigned a random name by RabbitMQ.
func SetQueueName(name string) {
	queueName = name
}

// declareQueue returns a queue. This function is idempotent, meaning
// it creates a queue if the queue does not already exist, or ensures
// that an existing queue matches the given config.
func declareQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
}
