package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
