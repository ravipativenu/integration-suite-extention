package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ravipativenu/integration-suite-extension/data"
	"ravipativenu/integration-suite-extension/jobs"
	"ravipativenu/integration-suite-extension/tests"
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

func testcases(w http.ResponseWriter, r *http.Request) {
	var t data.RawTestCase
	switch r.Method {
	case "GET":
		log.Println("Inside testcases GET handler...")
		w.Write(tests.GetTestCases())
	case "POST":
		log.Println("Inside testcases POST handler...")
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		tests.CreateTestCase(t)
		fmt.Fprintf(w, "TestCase: %+v", t)
	}

}

func scenarios(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside scenarios handler...")
	w.Write(tests.GetScenarios())
}

func payload(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside payload handler...")
	log.Println(r.URL.Query().Get("id"))
	w.Write(tests.GetPayload(r.URL.Query().Get("id")))
}

// Setting up the routes for incoming requests.
func setupRoutes() {
	http.HandleFunc("/", index)
	http.HandleFunc("/api/hello", hello)
	http.HandleFunc("/api/jobs/scheduler", jobscheduler)
	http.HandleFunc("/api/testing/testcases", testcases)
	http.HandleFunc("/api/testing/scenarios", scenarios)
	http.HandleFunc("/api/testing/payload", payload)
}

// The main method of the GO application setting up the routes and starting the http server.
func main() {
	fmt.Println("Go Web App Started on Port 8080")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
