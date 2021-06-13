package queue

import "github.com/streadway/amqp"

func getQueueName() string {
	return "image_queue"
}

// declareQueue returns a anonymous queue. It is idempotent, meaning it creates a queue if
// it does not already exist, or ensures that an existing queue matches the given config.
func declareQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		getQueueName(), // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
}
