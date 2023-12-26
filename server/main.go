package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/healthcheck", healthCheck)

	slog.Debug("Started server at http://localhost:8000")

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w)
}
