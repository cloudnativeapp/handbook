package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	username := os.Getenv("USERNAME")
	if username == "" {
		username = "world"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s\n", username)
	})
	http.ListenAndServe(":"+port, nil)
}
