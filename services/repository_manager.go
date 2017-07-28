package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type repositoryManager interface {
	retrievePodVersion() (string, error)
}

func RetrievePodVersion(rm repositoryManager) (string, error) {
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

func (n Nexus) retrievePodVersion() (string, error) {
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
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Nexus status code: %d", response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var nexusResponse NexusResponse
	err = json.Unmarshal(body, &nexusResponse)
	if err != nil {
		return "", err
	}

	if nexusResponse.Version == "" {
		return "", errors.New("Nexus: Version not found in the response")
	}

	return nexusResponse.Version, nil
}
