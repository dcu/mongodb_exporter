package shared

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Exportable defines an interface to export metrics to prometheus.
type Exportable interface {
	Describe(ch chan<- *prometheus.Desc)
	Export()
}
