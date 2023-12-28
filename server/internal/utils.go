package internal

import (
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func readTemplate(fileName string) string {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string
	return string(fileContent)
}

func connectToRabbitMQ(dsn string) *amqp.Connection {
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

	return conn
}
