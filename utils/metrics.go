package utils

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type MetricsRegistry struct {
	url      string
	job      string
	grouping map[string]string
	metrics  []prometheus.Collector
}

func NewBlackholeMetricsRegistry() *MetricsRegistry {
	return &MetricsRegistry{}
}

func NewPrometheusMetricsRegistry(config PrometheusConfig) *MetricsRegistry {
	return &MetricsRegistry{
		job:      "contactkey",
		grouping: push.HostnameGroupingKey(),
		url:      config.Url,
	}
}

func (mr *MetricsRegistry) Add(c prometheus.Collector) {
	mr.metrics = append(mr.metrics, c)
}

func (mr MetricsRegistry) Push() error {
	if mr.url == "" {
		return nil
	}
	return push.Collectors(mr.job, mr.grouping, mr.url, mr.metrics...)
}
