package queue

import (
	"log"

	"github.com/streadway/amqp"
)

type Producer struct {
	conn *amqp.Connection
}

func (p *Producer) start() error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()
	return nil
}

// Publish a message to the AMQP exchange.
func (p *Producer) Publish(msg string) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	q, err := declareQueue(ch)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(msg),
		},
	)
	if err != nil {
		return err
	}
	log.Printf("sending message: %s -> %s", msg, q.Name)

	return nil
}

func NewProducer(conn *amqp.Connection) (Producer, error) {
	producer := Producer{
		conn: conn,
	}
	if err := producer.start(); err != nil {
		return Producer{}, err
	}
	return producer, nil
}
