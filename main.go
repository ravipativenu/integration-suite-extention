package main

import (
	"fmt"
	"net/http"
	"ravipativenu/integration-suite-extension/jobs"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func jobscheduler(w http.ResponseWriter, r *http.Request) {
	jobs.RestartScheduler()
	fmt.Fprintf(w, "Scheduler Restarted...")
}

// Setting up the routes for incoming requests.
func setupRoutes() {
	http.HandleFunc("/", index)
	http.HandleFunc("/api/hello", hello)
	http.HandleFunc("/api/jobs/scheduler", jobscheduler)
}

// The main method of the GO application setting up the routes and starting the http server.
func main() {
	fmt.Println("Go Web App Started on Port 8080")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
