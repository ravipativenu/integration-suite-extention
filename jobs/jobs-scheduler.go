package jobs

import (
	"log"
	"ravipativenu/integration-suite-extension/data"
	"time"
)

var sc = make(chan string)

func RestartScheduler() {
	sc <- "restart"
}

func InitializeScheduler(r <-chan string) {
	//Get scheduler configuration from backend using GetSchedulerConfig()
	t, err := data.GetSchedulerConfig()
	if err != nil {
		log.Fatalln(err)
	}
	//Run scheduler in loop
	for {
		select {
		case <-r:
			//If any configuration update, Get scheduler configuration from backend
			//using GetSchedulerConfig()
			t, err = data.GetSchedulerConfig()
			if err != nil {
				log.Fatalln(err)
			}
		default:
			//Process sleep for scheduler time
			time.Sleep(t)
			//Schedule jobs
			scheduleJobs()
		}

	}
}

func scheduleJobs() {
	jobs, err := data.GetJobsList()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(jobs)
	go ScheduleJob(jobs[0])
}

//Schedules a job cheking some conditions
func ScheduleJob(job data.Job) {
	log.Printf("Scheduling Job %s", job.Name)
	log.Println(job)
	jobrun, err := data.GetLatestJobRun(job)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(jobrun)
	if jobrun.RunId == 0 {
		//Job is running for first time
		executeJob(job)
	} else {
		log.Printf("Previous scheduled run of job %s having RUNID %d started at %s", job.Name, jobrun.RunId, jobrun.StartTime)
		//Verify if the last job is not finished yet
		if jobrun.EndTime.IsZero() {
			log.Printf("Previous scheduled run of job %s with RUNID %d started at %s is still running", job.Name, jobrun.RunId, jobrun.StartTime)
			return
		}
		//Current time - Job schedule time. e.g. curent time - 30 min
		t := time.Now().Add(time.Minute * time.Duration(-1*job.Time))
		//Convert the database time from UTC to CST to compare with the current time of process
		loc, err := time.LoadLocation("America/Chicago")
		if err != nil {
			log.Fatal(err)
		}
		//Check if a job already ran within the scheduled time. e.g. within 10 mins. If yes, wait till scheduled time.
		if jobrun.StartTime.In(loc).After(t) {
			log.Printf("Job %s is scheduled to run every %d min. Need to wait for %f min for next schedueld run", job.Name, job.Time, (jobrun.StartTime.Sub(t)).Minutes())
			return
		}
		//Execute job is all conditions are satisfied
		executeJob(job)
	}
}

func executeJob(job data.Job) {
	switch job.Name {
	case "EXTRACTINTEGRATIONFLOWS":
		getIntegrationFlows(job)
	default:
		log.Printf("Job %s is not set up yet", job.Name)
	}
}

func init() {
	//go InitializeScheduler(sc)
}
