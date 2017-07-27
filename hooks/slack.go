package hooks

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Slack struct {
	Url     string
	Token   string
	Channel string
}

func (s Slack) preDeployment() error {
	// @todo Change the message
	return s.postMessage("Start update")
}

func (s Slack) postDeployment() error {
	// @todo Change the message
	return s.postMessage("End update")
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
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Slack status code: %d", response.StatusCode))
	}

	return nil
}
