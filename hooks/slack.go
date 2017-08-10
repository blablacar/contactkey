package hooks

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/remyLemeunier/contactkey/utils"
)

type Slack struct {
	Url     string
	Token   string
	Channel string
}

func NewSlack(cfg utils.SlackConfig, manifest utils.SlackManifest) *Slack {
	return &Slack{
		Url:     cfg.Url,
		Token:   cfg.Token,
		Channel: manifest.Channel,
	}
}

// @TODO Later we could pass directly messages and use the go templater instead.
func (s Slack) PreDeployment(username string, env string, service string, podVersion string) error {
	return s.postMessage(fmt.Sprintf("[%q]Start, update service %q version %q %q", env, service, podVersion, username))
}

func (s Slack) PostDeployment(username string, env string, service string, podVersion string) error {
	return s.postMessage(fmt.Sprintf("[%q]End, update service %q version %q by %q", env, service, podVersion, username))
}

func (s Slack) postMessage(message string) error {
	client := &http.Client{}
	url := fmt.Sprintf(
		"%s/services/hooks/slackbot?token=%s&channel=%s",
		s.Url,
		s.Token,
		s.Channel,
	)

	request, err := http.NewRequest("POST", url, strings.NewReader(message))
	if err != nil {
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Slack status code: %d", response.StatusCode))
	}

	return nil
}
