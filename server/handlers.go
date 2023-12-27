package main

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w)
}
