package internal

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
}

func NewConsumer(dsn string, queueName string) *Consumer {
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

	return &Consumer{conn, ch, queue}
}

func (c *Consumer) Consume(results chan []byte) {
	rawMsgs, err := c.ch.Consume(
		c.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range rawMsgs {
			log.Printf("received: %s", d.Body)
			results <- d.Body
		}
	}()

	log.Printf(" [*] Waiting for messages")
}

func (c *Consumer) Close() {
	c.ch.Close()
	c.conn.Close()
}
