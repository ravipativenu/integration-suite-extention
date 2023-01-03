package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"ravipativenu/integration-suite-extension/data"
	"ravipativenu/integration-suite-extension/jobs"
	"ravipativenu/integration-suite-extension/kafka"
	"ravipativenu/integration-suite-extension/tests"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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

func onboardSubscriber(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside subscribe handler...")
	params := mux.Vars(r)
	tenantId := params["TENANT_ID"]
	subscriberId := params["SUBSCRIBER_ID"]
	tenantInfo := data.TenantInfo{tenantId, subscriberId}
	response := data.Onboard(r.Context(), tenantId, tenantInfo)
	fmt.Fprintf(w, response)
}

func offboardSubscriber(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside unsubscribe handler...")
	params := mux.Vars(r)
	tenantId := params["TENANT_ID"]
	response := data.Offboard(r.Context(), tenantId)
	fmt.Fprintf(w, response)
}

// Setting up the routes for incoming requests.
func setupRoutes(router *mux.Router) {
	router.HandleFunc("/", index)
	router.HandleFunc("/api/hello", hello)
	router.HandleFunc("/api/jobs/scheduler", jobscheduler)
	router.HandleFunc("/api/testing/testcases", testcases)
	router.HandleFunc("/api/testing/scenarios", scenarios)
	router.HandleFunc("/api/testing/payload", payload)
	router.HandleFunc("/api/onboarding/rest/dbservice/v1/{TENANT_ID}/onboard/{SUBSCRIBER_ID}", onboardSubscriber).Methods("PUT")
	router.HandleFunc("/api/offboarding/rest/dbservice/v1/{TENANT_ID}/offboard", offboardSubscriber).Methods("PUT")
}

// The main method of the GO application setting up the routes and starting the http server.
func main() {

	var err error
	// load .env file
	err = godotenv.Load(".env")
	if err != nil {
		log.Println("No Local .env file. So accessing environment variables from Kyma runtime")
	}

	if goDotEnvVariable("KAFKA_ENABLE") == "true" {
		log.Println("Kafka is enabled thorugh environment...")
		ctx := context.Background()
		//go kafka.Produce2(ctx)
		go kafka.Consume2(ctx)
	}

	// Init the mux router
	router := mux.NewRouter()

	fmt.Println("Go Web App Started on Port 8080")
	setupRoutes(router)
	http.ListenAndServe(":8080", router)

}

func goDotEnvVariable(key string) string {
	return os.Getenv(key)
}
