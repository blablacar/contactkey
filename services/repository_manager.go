package services

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"strings"

	"github.com/remyLemeunier/contactkey/utils"
)

type RepositoryManager interface {
	RetrievePodVersion(sha1 string) (string, error)
}

type Nexus struct {
	Url        string
	Repository string
	Artifact   string
	Group      string
}

type NexusResponse struct {
	ArtifactId string `xml:"artifactId"`
	GroupId    string `xml:"groupId"`
	Versioning struct {
		Latest   string   `xml:"latest"`
		Release  string   `xml:"release"`
		Versions []string `xml:"versions>version"`
	} `xml:"versioning"`
}

func NewNexus(cfg utils.NexusConfig, manifest utils.NexusManifest) *Nexus {
	return &Nexus{
		Url:        cfg.Url,
		Repository: cfg.Repository,
		Group:      cfg.Group,
		Artifact:   manifest.Artifact,
	}
}

func (n Nexus) RetrievePodVersion(sha1 string) (string, error) {
	group := strings.Replace(n.Group, ".", "/", -1)
	client := &http.Client{}
	url := fmt.Sprintf(
		"%s/nexus/service/local/repositories/%s/content/%s/%s/maven-metadata.xml",
		n.Url,
		n.Repository,
		group,
		n.Artifact,
	)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

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
	err = xml.Unmarshal(body, &nexusResponse)
	if err != nil {
		return "", err
	}

	if sha1 != "" {
		for _, version := range nexusResponse.Versioning.Versions {
			if strings.Contains(version, sha1) {
				return version, nil
			}
		}

		return "", nil
	}

	return nexusResponse.Versioning.Latest, nil
}
