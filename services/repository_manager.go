package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/remyLemeunier/contactkey/utils"
)

type RepositoryManager interface {
	RetrievePodVersion() (string, error)
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

func NewNexus(cfg utils.NexusConfig, manifest utils.NexusManifest) *Nexus {
	return &Nexus{
		Url:        cfg.Url,
		Repository: cfg.Repository,
		Group:      cfg.Group,
		Artifact:   manifest.Artifact,
	}
}

func (n Nexus) RetrievePodVersion() (string, error) {
	client := &http.Client{}
	url := fmt.Sprintf(
		"%s/nexus/service/local/artifact/maven?r=%s&a=%s&g=%s&v=LATEST",
		n.Url,
		n.Repository,
		n.Artifact,
		n.Group,
	)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

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
