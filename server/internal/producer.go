package internal

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
}

const MAX_RETRIES = 5

func NewProducer(dsn string, queueName string) *Producer {
	// opening connection
	conn := connectToRabbitMQ(dsn)

	// opening channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	queue, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		true,      // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	return &Producer{conn, ch, queue}
}

func (p *Producer) Send(msg []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.ch.PublishWithContext(ctx,
		"",           // exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})
	if err != nil {
		return err
	}

	log.Printf(" [x] sent %s\n", msg)

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
