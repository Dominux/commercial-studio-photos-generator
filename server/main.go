package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// starting producer
	amqp_dsn := load_amqp_dsn()
	qName := os.Getenv("RABBITMQ_QUEUE")
	producer := NewProducer(amqp_dsn, qName)
	defer producer.Close()

	http.HandleFunc("/", index)
	http.HandleFunc("/healthcheck", healthCheck)
	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "404 not Found", http.StatusNotFound)
			return
		}

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		product := r.FormValue("product")
		id, err := runGenerating(product, producer)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}

		cookie := &http.Cookie{Name: "id", Value: id.String()}
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusNoContent)
	})

	log.Print("Started server at http://localhost:8000")

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func load_amqp_dsn() string {
	amqp_host := os.Getenv("RABBITMQ_HOST")
	amqp_port := os.Getenv("RABBITMQ_PORT")
	amqp_user := os.Getenv("RABBITMQ_USER")
	amqp_pass := os.Getenv("RABBITMQ_PASSWORD")

	return fmt.Sprintf("amqp://%s:%s@%s:%s/", amqp_user, amqp_pass, amqp_host, amqp_port)
}
