package queue

import (
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn *amqp.Connection
}

func (c *Consumer) start() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()
	return nil
}

// Listen starts listening for messages on the queue and runs the given job
// for each each new message received.
func (c *Consumer) Listen(job func(d amqp.Delivery) error) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareQueue(ch)
	if err != nil {
		return err
	}

	// Fairly dispatch messages to all workers as they are not busy.
	ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			if err := job(d); err != nil {
				log.Printf("error processing %d: %s", d.DeliveryTag, err)
				d.Reject(false)
			}
			d.Ack(false)
		}
	}()

	log.Printf("[*] Waiting for message [Queue][%s]. To exit press CTRL+C", q.Name)
	<-forever
	return nil
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}
	if err := consumer.start(); err != nil {
		return Consumer{}, err
	}
	return consumer, nil
}
