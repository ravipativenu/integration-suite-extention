package data

import (
	"fmt"
	"log"
	"time"
)

func ControlJob(cmd string, rid int, jobid int) (id int) {
	var stim time.Time
	db, err := db.getDb()
	if err != nil {
		log.Fatal(err)
	}
	switch cmd {
	case "START":
		if err := db.QueryRow(fmt.Sprintf("SELECT %s.NEXTVAL AS ID FROM DUMMY", "INTEGRATION_SUITE.JOBRUNID")).Scan(&rid); err != nil {
			log.Fatal(err)
		}
		stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s values (?, ?, ?, ?, ?, ?, ?)", "INTEGRATION_SUITE.JOBRUNS"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		if _, err := stmt.Exec(rid, jobid, time.Now(), nil, nil, "RUNNING", "/logurl"); err != nil {
			log.Fatal(err)
		}
		return rid
	case "END":
		err := db.QueryRow(fmt.Sprintf("SELECT %s FROM %s WHERE RUNID = %d", "starttime", "INTEGRATION_SUITE.JOBRUNS", rid)).Scan(&stim)
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := db.Prepare(fmt.Sprintf("UPDATE %s SET endtime = ?, duration = ?, STATUS = ? WHERE RUNID = ?", "INTEGRATION_SUITE.JOBRUNS"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		log.Println(time.Now().Sub(stim).String())
		if _, err := stmt.Exec(time.Now(), time.Now().Sub(stim).String(), "FINISHED", rid); err != nil {
			log.Fatal(err)
		}
		return rid
	case "ERROR":
		err := db.QueryRow(fmt.Sprintf("SELECT %s FROM %s WHERE RUNID = %d", "starttime", "INTEGRATION_SUITE.JOBRUNS", rid)).Scan(&stim)
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := db.Prepare(fmt.Sprintf("UPDATE %s SET endtime = ?, duration = ?, STATUS = ? WHERE RUNID = ?", "INTEGRATION_SUITE.JOBRUNS"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		log.Println(time.Now().Sub(stim).String())
		if _, err := stmt.Exec(time.Now(), time.Now().Sub(stim).String(), "ERROR", rid); err != nil {
			log.Fatal(err)
		}
		return rid
	default:
		log.Fatal("Only START and END are supported")
	}
	return rid
}
