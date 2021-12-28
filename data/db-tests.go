package data

import (
	"fmt"
	"log"
)

type RawTestCase struct {
	Name        string
	Description string
	Filedata    string
	Filesize    int
	Filename    string
	Filetype    string
	Method      string
	Testcase    string
}

type Testcase struct {
	ID          int
	Name        string
	Testcase    string
	Description string
	Method      string
	Input       string
}

type Scenario struct {
	ID          int
	Name        string
	Version     string
	Description string
	Protocol    string
	Endpoint    string
}

func GetTestCases() ([]Testcase, error) {
	db, err := db.getDb()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", "INTEGRATION_SUITE.TESTCASES"))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	//A testcase slice to hold data from returned rows
	var testcases []Testcase
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var testcase Testcase
		if err := rows.Scan(&testcase.ID, &testcase.Name, &testcase.Testcase, &testcase.Description, &testcase.Method, &testcase.Input); err != nil {
			return testcases, err
		}
		testcases = append(testcases, testcase)
	}
	if err = rows.Err(); err != nil {
		return testcases, err
	}
	return testcases, nil
}

func GetScenarios() ([]Scenario, error) {
	db, err := db.getDb()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", "INTEGRATION_SUITE.IFLOWS"))
	log.Println(rows)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	//A testcase slice to hold data from returned rows
	var scenarios []Scenario

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var scenario Scenario
		if err := rows.Scan(&scenario.ID, &scenario.Name, &scenario.Version, &scenario.Description, &scenario.Protocol, &scenario.Endpoint); err != nil {
			return scenarios, err
		}
		scenarios = append(scenarios, scenario)
	}
	if err = rows.Err(); err != nil {
		return scenarios, err
	}
	return scenarios, nil
}

func CreateTestCase(t RawTestCase) {
	db, err := db.getDb()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s (NAME, TESTCASE, DESCRIPTION, METHOD, INPUT) values (?, ?, ?, ?, ?)", "INTEGRATION_SUITE.TESTCASES"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(t.Name, t.Testcase, t.Description, t.Method, "/"+t.Name+"/"+t.Testcase+"/"+t.Filename); err != nil {
		log.Fatal(err)
	}
}
