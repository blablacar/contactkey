package hooks

import (
	"errors"
	"fmt"
	"github.com/golang/glog"
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
	if err != nil {
		return err
	}

	glog.V(4).Infof("sending %q to %s", message, url)
	response, err := client.Do(request)
	if err != nil {
		glog.Errorf("fail to send %q to %s: %q", message, url, err)
		return err
	}

	if response.StatusCode != http.StatusOK {
		glog.Errorf("fail to send %q to %s StatusCode: %d", message, url, response.StatusCode)
		return errors.New(fmt.Sprintf("Slack status code: %d", response.StatusCode))
	}

	return nil
}
