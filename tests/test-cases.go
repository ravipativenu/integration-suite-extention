package tests

import (
	"encoding/json"
	"log"
	"ravipativenu/integration-suite-extension/data"
)

type RawTestCase struct {
	Name        string
	Description string
	Filedata    string
	Filesize    int
	Filetype    string
	Method      string
	Testcase    string
}

func GetTestCases() []byte {
	testcases, err := data.GetTestCases()
	if err != nil {
		log.Fatalln(err)
	}
	data, err := json.Marshal(testcases)
	if err != nil {
		log.Fatalln(err)
	}
	return (data)
}

func GetScenarios() []byte {
	log.Println("tests.GetScenarios...")
	scenarios, err := data.GetScenarios()
	if err != nil {
		log.Fatalln(err)
	}
	data, err := json.Marshal(scenarios)
	if err != nil {
		log.Fatalln(err)
	}
	return (data)
}

func CreateTestCase(t data.RawTestCase) {
	data.CreateTestCase(t)
	data.CreateTestCaseBlob(t)
}

func GetPayload(f string) []byte {
	payload := data.GetTestCasePayloadBlob(f)
	return payload
}
