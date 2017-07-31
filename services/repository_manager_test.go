package services

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"io/ioutil"
)

func TestNexus(t *testing.T) {
	var apiStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/nexus/service/local/artifact/maven?r=repository&a=artifact&g=group&v=LATEST":
			raw, err := ioutil.ReadFile("./data/result.json")
			if err != nil {
				t.Error(err)
			}

			w.WriteHeader(http.StatusOK)
			w.Write(raw)
		}
	}))

	nexus := Nexus{
		apiStub.URL,
		"repository",
		"artifact",
		"group",
	}

	podVersion, err := nexus.retrievePodVersion()
	if err != nil {
		t.Errorf("Error triggered %s", err.Error())
	}

	if podVersion != "26.1501244191-vb0f586a" {
		t.Errorf("podVersion is different from expected podVersion: %s", podVersion)
	}

}
