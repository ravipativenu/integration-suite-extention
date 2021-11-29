package data

import (
	"database/sql"
	"log"
	"os"

	// Register hdb driver.
	_ "github.com/SAP/go-hdb/driver"
	"github.com/joho/godotenv"
)

type HanaDBEnv struct {
	driverName string
	hdbDsn     string
	Db         *sql.DB
}

type DbSource interface {
	getDb() (*sql.DB, error)
}

var db = &HanaDBEnv{"", "", nil}

func (h *HanaDBEnv) getDb() (*sql.DB, error) {
	if h.Db == nil && h.driverName == "" && h.hdbDsn == "" {
		var err error
		// load .env file
		err = godotenv.Load(".env")
		if err != nil {
			log.Println("No Local .env file. So accessing environment variables from Kyma runtime")
		}
		h.driverName = goDotEnvVariable("HANA_SECRET_DRIVER")
		h.hdbDsn = goDotEnvVariable("HANA_SECRET_DSN")
		Db, err := sql.Open(h.driverName, h.hdbDsn)
		if err != nil {
			return nil, err
		}
		h.Db = Db
		return h.Db, nil
	} else {
		return h.Db, nil
	}

}

func goDotEnvVariable(key string) string {
	return os.Getenv(key)
}
