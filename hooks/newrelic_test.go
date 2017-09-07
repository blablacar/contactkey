package hooks

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var GetApplications = []byte(`{
  "applications": [
    {
      "id": 456,
      "name": "Webhooks"
    }
  ]
}`)

var GetApplicationsEmpty = []byte(`{"applications": []}`)

var apiStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	if r.Header.Get("X-Api-Key") != "123" {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		switch r.RequestURI {
		case "/v2/applications.json":
			w.WriteHeader(http.StatusOK)
			if string(body) == "filter[name]=webhook" {
				w.Write(GetApplications)
			} else {
				w.Write(GetApplicationsEmpty)
			}

		case "/v2/applications/456/deployments.json":
			w.WriteHeader(http.StatusCreated)

		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}))

func TestCreateDeployment(t *testing.T) {
	c := &NewRelicClient{
		HttpClient:    &http.Client{},
		Url:           apiStub.URL,
		ApplicationId: 456,
		ApiKey:        "123",
		Stop:          false,
	}
	d := &NewRelicDeployment{}

	err := c.CreateDeployment(d)
	if err != nil {
		t.Errorf("Unexpected err : %q", err)
	}
}

func TestFindApplicationId(t *testing.T) {
	c := &NewRelicClient{
		HttpClient: &http.Client{},
		Url:        apiStub.URL,
		ApiKey:     "123",
		Stop:       false,
	}
	// log.SetLevel(log.DebugLevel)

	appId, err := c.findApplicationId("webhook")
	if err != nil {
		t.Errorf("Unexpected err : %q", err)
	}
	if appId != 456 {
		t.Errorf("Unexpected appId : %q", appId)
	}
}
