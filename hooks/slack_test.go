package hooks

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSlack(t *testing.T) {
	var apiStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/services/hooks/slackbot?token=abc&channel=channel":
			w.WriteHeader(http.StatusOK)

		case "/services/hooks/slackbot?token=abc&channel=error400":
			w.WriteHeader(http.StatusBadRequest)
		}
	}))

	slack := &Slack{
		apiStub.URL,
		"abc",
		"channel",
		false,
	}
	err := slack.postMessage("Some message.")
	if err != nil {
		t.Error("Shouldn't have recieved an err")
	}

	// Chance the channel to trigger an http.statusBadRequest
	slack.Channel = "error400"
	err = slack.postMessage("Some message.")
	if err.Error() != "Slack status code: 400" {
		t.Error("It should have triggered an error due to code 400")
	}
}
