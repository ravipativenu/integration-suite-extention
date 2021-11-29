package data

import (
	"fmt"
	"log"
)

type IFlow struct {
	Name        string
	Version     string
	Description string
	Protocol    string
	Endpoint    string
}

func UpdateIFlows(iflows []IFlow) {
	db, err := db.getDb()
	if err != nil {
		log.Fatal(err)
	}
	// Truncate table.
	if _, err := db.Exec(fmt.Sprintf("TRUNCATE table %s", "INTEGRATION_SUITE.IFLOWS")); err != nil {
		log.Fatal(err)
	}
	log.Println(iflows)
	stmt, err := db.Prepare(fmt.Sprintf("BULK INSERT INTO %s (Name, Version, Description, Protocol, Endpoint) values (?, ?, ?, ?, ?)", "INTEGRATION_SUITE.IFLOWS"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	// Bulk insert.
	for i := 0; i < len(iflows); i++ {
		if _, err := stmt.Exec(iflows[i].Name, iflows[i].Version, iflows[i].Description, iflows[i].Protocol, iflows[i].Endpoint); err != nil {
			log.Fatal(err)
		}
	}
	// Call final stmt.Exec().
	if _, err := stmt.Exec(); err != nil {
		log.Fatal(err)
	}
	log.Printf("%d rows inserted in to table %s", len(iflows), "INTEGRATION_SUITE.IFLOWS")
}
