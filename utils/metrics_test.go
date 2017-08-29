package utils

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

var apiStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if m, _ := regexp.MatchString(`/metrics/job/contactkey/instance/.*`, r.RequestURI); m == true {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}))

func TestMetricsPush(t *testing.T) {
	mr := NewPrometheusMetricsRegistry(
		PrometheusConfig{
			Url: apiStub.URL,
		})
	mr.Push()
}
