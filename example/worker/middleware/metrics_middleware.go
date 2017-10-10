package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tkanos/konsumerou"
)

const (
	requestName       = "request_total"
	requestFailedName = "request_failed"
	latencyName       = "request_latency_milliseconds"
)

type metricsService struct {
	request       *prometheus.CounterVec
	requestFailed *prometheus.CounterVec
	latency       *prometheus.SummaryVec
	serviceName   string
}

// NewMetricsService creates a layer of service that add metrics capability
func NewMetricsService(serviceName string, next konsumerou.Handler) konsumerou.Handler {
	m := metricsMiddleware(serviceName)
	return m.instrumentation(next)
}

func metricsMiddleware(name string) *metricsService {
	var m metricsService
	fieldKeys := []string{"method"}

	m.request = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "Consumer",
			Subsystem: "my_app",
			Name:      fmt.Sprintf("%v_%v", strings.Replace(name, "-", "_", -1), requestName),
			Help:      "Number of requests processed",
		}, fieldKeys)
	prometheus.MustRegister(m.request)

	m.requestFailed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "Consumer",
			Subsystem: "my_app",
			Name:      fmt.Sprintf("%v_%v", strings.Replace(name, "-", "_", -1), requestFailedName),
			Help:      "Number of requests failed",
		}, fieldKeys)
	prometheus.MustRegister(m.requestFailed)

	m.latency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "Consumer",
			Subsystem: "my_app",
			Name:      fmt.Sprintf("%v_%v", strings.Replace(name, "-", "_", -1), latencyName),
			Help:      "Total duration in miliseconds.",
		}, fieldKeys)
	prometheus.MustRegister(m.latency)

	m.serviceName = name

	return &m
}

func (m *metricsService) instrumentation(next konsumerou.Handler) konsumerou.Handler {
	return func(ctx context.Context, msg *sarama.ConsumerMessage) (err error) {
		start := time.Now()
		// add metrics to this method
		defer m.latency.WithLabelValues(m.serviceName).Observe(time.Since(start).Seconds() * 1e3)
		defer m.request.WithLabelValues(m.serviceName).Inc()
		// If error is not empty, we add to metrics that it failed
		if err = next(ctx, msg); err != nil {
			m.requestFailed.WithLabelValues(m.serviceName).Inc()
		}
		return
	}
}
