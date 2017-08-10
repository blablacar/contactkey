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
		case "/nexus/service/local/repositories/repo/content/this/is/a/group/service-name/maven-metadata.xml":
			raw, err := ioutil.ReadFile("./data/artifact_metadata.xml")
			if err != nil {
				t.Error(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(raw)
		}
	}))

	nexus := Nexus{
		apiStub.URL,
		"repo",
		"service-name",
		"this.is.a.group",
	}

	podVersion, err := nexus.RetrievePodVersion("")
	if err != nil {
		t.Fatalf("Error triggered %s", err.Error())
	}

	if podVersion != "26.1502293426-vb2d95f2" {
		t.Errorf("podVersion is different from expected podVersion: %s", podVersion)
	}

	podVersion, err = nexus.RetrievePodVersion("60eccd7")
	if err != nil {
		t.Fatalf("Error triggered %s", err.Error())
	}

	if podVersion != "26.1499694197-v60eccd7" {
		t.Errorf("podVersion is different from expected podVersion: %s", podVersion)
	}

	podVersion, err = nexus.RetrievePodVersion("toto")
	if err != nil {
		t.Fatalf("Error triggered %s", err.Error())
	}

	if podVersion != "" {
		t.Error("podVersion should be empty")
	}
}
