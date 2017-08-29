package hooks

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var apiStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.RequestURI {
	case "/v2/applications/456/deployments.json":
		if r.Header.Get("X-Api-Key") != "123" {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}))

func TestCreateDeployment(t *testing.T) {
	c := &NewRelicClient{
		Url:           apiStub.URL,
		ApiKey:        "123",
		ApplicationId: "456",
		Stop:          false,
	}
	//c.Log.SetLevel(log.DebugLevel)
	d := &NewRelicDeployment{}

	err := c.CreateDeployment(d)
	if err != nil {
		t.Errorf("Unexpected err : %q", err)
	}
}
