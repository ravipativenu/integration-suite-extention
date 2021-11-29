package jobs

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"ravipativenu/integration-suite-extension/data"

	"golang.org/x/oauth2/clientcredentials"
)

type ServiceEndpointsResponse struct {
	D struct {
		Results []struct {
			Metadata struct {
				ID   string `json:"id"`
				URI  string `json:"uri"`
				Type string `json:"type"`
			} `json:"__metadata"`
			Name        string `json:"Name"`
			ID          string `json:"Id"`
			Title       string `json:"Title"`
			Version     string `json:"Version"`
			Summary     string `json:"Summary"`
			Description string `json:"Description"`
			LastUpdated string `json:"LastUpdated"`
			Protocol    string `json:"Protocol"`
			EntryPoints struct {
				Results []struct {
					Metadata struct {
						ID   string `json:"id"`
						URI  string `json:"uri"`
						Type string `json:"type"`
					} `json:"__metadata"`
					Name                  string `json:"Name"`
					URL                   string `json:"Url"`
					Type                  string `json:"Type"`
					AdditionalInformation string `json:"AdditionalInformation"`
				} `json:"results"`
			} `json:"EntryPoints"`
			APIDefinitions struct {
				Deferred struct {
					URI string `json:"uri"`
				} `json:"__deferred"`
			} `json:"ApiDefinitions"`
		} `json:"results"`
	} `json:"d"`
}

func getIntegrationFlows(job data.Job) {
	log.Printf("Job %s Started", job.Name)
	//When job started, update backend the job started details and get the runid assinged for the job run.
	rid := data.ControlJob("START", -1, job.ID)
	ctx := context.Background()
	conf := clientcredentials.Config{
		ClientID:     os.Getenv("CPI_SECRET_CLIENTID"),
		ClientSecret: os.Getenv("CPI_SECRET_CLIENTSECRET"),
		TokenURL:     os.Getenv("CPI_SECRET_TOKENENDPOINT"),
	}
	client := conf.Client(ctx)
	resp, err := client.Get(os.Getenv("CPI_SECRET_APIENDPOINT") + "/api/v1/ServiceEndpoints?$expand=EntryPoints&$format=json")
	if err != nil {
		data.ControlJob("ERROR", rid, job.ID)
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	var result ServiceEndpointsResponse
	if err = json.Unmarshal(body, &result); err != nil {
		data.ControlJob("ERROR", rid, job.ID)
		log.Fatalln(err)
	}
	var iflows []data.IFlow
	// Loop through the api response for the FirstName
	for _, rec := range result.D.Results {
		var iflow data.IFlow
		iflow.Name = rec.Name
		iflow.Version = rec.Version
		iflow.Description = rec.Description
		iflow.Protocol = rec.Protocol
		iflow.Endpoint = rec.EntryPoints.Results[0].URL
		iflows = append(iflows, iflow)
	}
	data.UpdateIFlows(iflows)
	//When job completed, update backend the job completion details.
	log.Printf("Job %s Finished", job.Name)
	data.ControlJob("END", rid, job.ID)
}
