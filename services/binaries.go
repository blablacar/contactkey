package services

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/remyLemeunier/contactkey/utils"
	log "github.com/sirupsen/logrus"
)

type Binaries interface {
	RetrievePodVersion(sha1 string) (string, error)
}

type Nexus struct {
	Url        string
	Repository string
	Artifact   string
	Group      string
	Log        *log.Logger
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

func NewNexus(cfg utils.NexusConfig, manifest utils.NexusManifest, logger *log.Logger) (*Nexus, error) {
	if cfg.Url == "" {
		return nil, errors.New("You need to define an url for nexus in the config.")
	}

	if cfg.Repository == "" {
		return nil, errors.New("You need to define a repository for nexus in the config.")
	}

	if cfg.Group == "" {
		return nil, errors.New("You need to define a group for nexus in the config.")
	}

	if manifest.Artifact == "" {
		return nil, errors.New("You need to define an artifact for nexus in the manifest.")
	}

	return &Nexus{
		Url:        cfg.Url,
		Repository: cfg.Repository,
		Group:      cfg.Group,
		Artifact:   manifest.Artifact,
		Log:        logger,
	}, nil
}

// If no sha1 is given to this function
// It will retrieve the LATEST version available.
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

	n.Log.WithFields(log.Fields{
		"fullPath":   url,
		"nexusUrl":   n.Url,
		"repository": n.Repository,
		"Group":      n.Group,
		"Artifact":   n.Artifact,
	}).Debug("Creating Nexus url")

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	n.Log.WithFields(log.Fields{
		"statusCode": response.StatusCode,
	}).Debug("Status code from Nexus")

	if response.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Nexus status code: %d", response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	n.Log.WithFields(log.Fields{
		"body": string(body),
	}).Debug("Response body from Nexus")

	var nexusResponse NexusResponse
	err = xml.Unmarshal(body, &nexusResponse)
	if err != nil {
		return "", err
	}

	serviceVersion := ""
	if sha1 != "" {
		for _, version := range nexusResponse.Versioning.Versions {
			if strings.Contains(version, sha1) {
				serviceVersion = version
				break
			}
		}
	} else {
		serviceVersion = nexusResponse.Versioning.Latest
	}

	n.Log.WithFields(log.Fields{
		"serviceVersion": serviceVersion,
	}).Debug("Version retrieved in Nexus")

	return serviceVersion, nil
}
