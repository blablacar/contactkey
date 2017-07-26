package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/gommon/log"
)

type repositoryManager interface {
	retrievePodVersion() string
}

func RetrievePodVersion(rm repositoryManager) string {
	return rm.retrievePodVersion()
}

type Nexus struct {
	Url        string
	Repository string
	Artifact   string
	Group      string
}

type NexusResponse struct {
	Version string `json:"version"`
}

func (n Nexus) retrievePodVersion() string {
	client := &http.Client{}
	url := fmt.Sprintf(
		"%s/nexus/service/local/artifact/maven?r=%s&a=%s&g=%s&v=LATEST",
		n.Url,
		n.Repository,
		n.Artifact,
		n.Group,
	)

	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("accept", "application/json")
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	var nexusResponse NexusResponse
	err = json.Unmarshal(body, &nexusResponse)
	if err != nil {
		log.Fatal(err.Error())
	}

	if nexusResponse.Version == "" {
		log.Fatal("Nexus: Version not found in the response")
	}

	return nexusResponse.Version
}
