package hooks

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/gommon/log"
)

type Slack struct {
	Url     string
	Token   string
	Channel string
}

func (s Slack) preDeployment() {
	// @todo Change the message
	s.postMessage("Start update")
}

func (s Slack) postDeployment() {
	// @todo change the message
	s.postMessage("End update")
}

func (s Slack) postMessage(message string) {
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
		log.Fatal(err.Error())
	}

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Slack status code: %s", http.StatusOK)
	}
}
