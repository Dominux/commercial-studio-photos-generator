package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	conn      *amqp.Connection
	ch        *amqp.Channel
	queueName string
}

const MAX_RETRIES = 5

func NewProducer(dsn string, queueName string) *Producer {
	// opening connection
	var (
		conn *amqp.Connection
		err  error
	)
	for i := 0; i < MAX_RETRIES; i++ {
		conn, err = amqp.Dial(dsn)
		if err == nil {
			break
		}

		time.Sleep(3 * time.Second)

		if i+1 == MAX_RETRIES {
			failOnError(err, "Failed to connect to RabbitMQ")
		}
	}

	// opening channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return &Producer{conn, ch, queueName}
}

func (p *Producer) Send(msg []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.ch.PublishWithContext(ctx,
		"",          // exchange
		p.queueName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})
	if err != nil {
		return err
	}

	log.Printf(" [x] Sent %s\n", msg)

	return nil
}

func (p *Producer) Close() {
	p.ch.Close()
	p.conn.Close()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
