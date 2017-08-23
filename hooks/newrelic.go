package hooks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/remyLemeunier/contactkey/utils"
	log "github.com/sirupsen/logrus"
)

type NewRelicClient struct {
	Url           string
	ApiKey        string
	ApplicationId string
	Log           *log.Logger
	Stop          bool
}

type NewRelicDeployment struct {
	revision    string
	changelog   string
	description string
	user        string
}

func (c NewRelicClient) PreDeployment(userName string, env string, service string, podVersion string) error {
	description := fmt.Sprintf("Deploying %s %s on %s", service, podVersion, env)
	d := &NewRelicDeployment{
		description: description,
		user:        userName,
	}
	return c.CreateDeployment(d)
}

func (c NewRelicClient) PostDeployment(userName string, env string, service string, podVersion string) error {
	return nil
}

func (c NewRelicClient) StopOnError() bool {
	return c.Stop
}

func NewNewRelicClient(cfg utils.NewRelicConfig, manifest utils.NewRelicManifest, logger *log.Logger) *NewRelicClient {
	return &NewRelicClient{
		Url:           cfg.Url,
		ApiKey:        cfg.ApiKey,
		ApplicationId: manifest.ApplicationId,
		Log:           logger,
		Stop:          manifest.StopOnError,
	}
}

// https://rpm.newrelic.com/api/explore/application_deployments/create
func (c NewRelicClient) CreateDeployment(d *NewRelicDeployment) error {
	client := &http.Client{}
	url := fmt.Sprintf("%s/v2/applications/%s/deployments.json",
		c.Url,
		c.ApplicationId,
	)

	c.Log.WithFields(log.Fields{
		"url": url,
	}).Debug("Creating NewRelic deployment.")

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(d); err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, nil)
	request.Header.Add("X-Api-Key", c.ApiKey)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	c.Log.WithFields(log.Fields{
		"statusCode": response.StatusCode,
	}).Debug("NewRelic response status code.")

	if response.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("NewRelic status code: %d", response.StatusCode))
	}

	return nil
}
