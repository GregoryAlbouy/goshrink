package main

import (
	"log"

	"github.com/GregoryAlbouy/shrinker/internal"
	"github.com/streadway/amqp"
)

type queueHandler struct {
	userService internal.UserService
}

func (h queueHandler) handleMessage(d amqp.Delivery) error {
	log.Println("Got avatar from " + d.MessageId)

	return nil
}
