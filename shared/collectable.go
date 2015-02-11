package shared

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Collectable interface {
	Collect(groupName string, ch chan<- prometheus.Metric)
}
