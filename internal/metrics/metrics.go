package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	UpdateEventLabel = "update"
	CreateEventLabel = "create"
	RemoveEventLabel = "remove"
)

var CounterCollector *prometheus.CounterVec

func CreateMetrics() {
	CounterCollector = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "reminds",
		Name:      "processed_events_total",
		Help:      "Counter with three event types",
	}, []string{UpdateEventLabel, CreateEventLabel, RemoveEventLabel})

	var registerer = prometheus.DefaultRegisterer
	registerer.MustRegister(CounterCollector)
}
