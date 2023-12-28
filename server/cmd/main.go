package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"

	"app/internal"
)

func main() {
	var amqp_dsn string
	{
		amqp_host := os.Getenv("RABBITMQ_HOST")
		amqp_port := os.Getenv("RABBITMQ_PORT")
		amqp_user := os.Getenv("RABBITMQ_USER")
		amqp_pass := os.Getenv("RABBITMQ_PASSWORD")

		amqp_dsn = fmt.Sprintf("amqp://%s:%s@%s:%s/", amqp_user, amqp_pass, amqp_host, amqp_port)
	}

	// setting redis client
	var redisClient *redis.Client
	{
		host := os.Getenv("REDIS_HOST")
		port := os.Getenv("REDIS_PORT")
		pass := os.Getenv("REDIS_PASS")
		addr := fmt.Sprintf("%s:%s", host, port)

		redisClient = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pass,
			DB:       0, // use default DB
		})
	}

	results := make(chan []byte)
	go internal.SaveResults(results, redisClient)

	// initing minio client
	var minioClient *minio.Client
	var err error
	{
		host := os.Getenv("S3_HOST")
		port := os.Getenv("S3_PORT")
		accessKeyID := os.Getenv("S3_ACCESS_KEY_ID")
		secretAccessKey := os.Getenv("S3_SECRET_ACCESS_KEY")
		endpoint := fmt.Sprintf("%s:%s", host, port)

		minioClient, err = minio.New(endpoint, &minio.Options{
			Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		})
		if err != nil {
			log.Fatalln(err)
		}
	}
	minioBucket := os.Getenv("S3_BUCKET")
	minioResultsPath := os.Getenv("S3_OUTPUTS_PATH")

	// starting consumer
	responseQName := os.Getenv("RABBITMQ_RESPONSE_QUEUE")
	consumer := internal.NewConsumer(amqp_dsn, responseQName)
	defer consumer.Close()
	consumer.Consume(results)

	// starting producer
	requestQName := os.Getenv("RABBITMQ_REQUEST_QUEUE")
	producer := internal.NewProducer(amqp_dsn, requestQName)
	defer producer.Close()

	// declaring handlers
	handler := internal.NewHandler(producer, redisClient, minioClient, minioBucket, minioResultsPath)
	{
		http.HandleFunc("/", handler.Index)
		http.HandleFunc("/healthcheck", handler.HealthCheck)
		http.HandleFunc("/generate", handler.RunGenerating)
		http.HandleFunc(internal.CheckProgressEndpoint, handler.CheckProgressHandler)
		http.HandleFunc(internal.GetResultEndpoint, handler.GetResultHandler)
	}

	// starting server
	log.Print("started server at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
