package main

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

// Setting up the routes for incoming requests.
func setupRoutes() {
	http.HandleFunc("/", index)
	http.HandleFunc("/api/hello", hello)
}

// The main method of the GO application setting up the routes and starting the http server.
func main() {
	fmt.Println("Go Web App Started on Port 8080")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
