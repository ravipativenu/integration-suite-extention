package data

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Job struct {
	ID       int
	Name     string
	Time     int
	Duration string
	Params   string
}

type JobRun struct {
	RunId     int
	JobId     int
	StartTime time.Time
	EndTime   time.Time
	Duration  string
	Status    string
	logUrl    string
}

//Function to get job scheduler configuration from backend or error if any failure
func GetSchedulerConfig() (time.Duration, error) {
	//Initialize time duration variable in nano seconds
	var t time.Duration = 1000000000
	db, err := db.getDb()
	if err != nil {
		return t, err
	}
	db.QueryRow(fmt.Sprintf("SELECT TIME FROM %s", "INTEGRATION_SUITE.JOBSCHEDULER")).Scan(&t)
	log.Printf("Scheduler configuration: %s", t)
	return t * 1000000000, nil
}

//Function to get jobs list from backend or error if any failure
func GetJobsList() ([]Job, error) {
	db, err := db.getDb()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", "INTEGRATION_SUITE.JOBSCHEDULES"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//A job slice to hold data from returd rows
	var jobs []Job

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var job Job
		if err := rows.Scan(&job.ID, &job.Name, &job.Time, &job.Duration, &job.Params); err != nil {
			return jobs, err
		}
		jobs = append(jobs, job)
	}
	if err = rows.Err(); err != nil {
		return jobs, err
	}
	return jobs, nil
}

//Get latest execution details of a job
func GetLatestJobRun(job Job) (JobRun, error) {
	var jobrun JobRun
	db, err := db.getDb()
	if err != nil {
		return jobrun, err
	}

	row := db.QueryRow(fmt.Sprintf("SELECT TOP 1 * FROM %s WHERE JOBID =%d ORDER BY STARTTIME DESC", "INTEGRATION_SUITE.JOBRUNS", job.ID))
	if row.Err() != nil {
		return jobrun, err
	}
	//Declare variable to take care of nullable values from database
	var etime sql.NullTime
	var dur sql.NullString
	err = row.Scan(&jobrun.RunId, &jobrun.JobId, &jobrun.StartTime, &etime, &dur, &jobrun.Status, &jobrun.logUrl)
	jobrun.EndTime = etime.Time
	jobrun.Duration = dur.String
	if err != nil {
		return jobrun, nil
	}
	return jobrun, nil
}
