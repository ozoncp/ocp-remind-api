package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	create,
	update,
	remove prometheus.Counter
)

func CreateMetrics() {
	create = promauto.NewCounter(prometheus.CounterOpts{Name: "create_reminds_counter"})
	update = promauto.NewCounter(prometheus.CounterOpts{Name: "update_reminds_counter"})
	remove = promauto.NewCounter(prometheus.CounterOpts{Name: "remove_reminds_counter"})
}

func CreateCounterUp() {
	create.Inc()
}

func UpdateCounterUp() {
	update.Inc()
}

func RemoveCounterUp() {
	remove.Inc()
}
