package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	cursorsMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "cursors",
		Help:      "The cursors data structure contains data regarding cursor state and use",
	}, []string{})
)

// Cursors are the cursor metrics
type Cursors struct {
	TotalOpen      float64 `bson:"totalOpen"`
	TimeOut        float64 `bson:"timedOut"`
	TotalNoTimeout float64 `bson:"totalNoTimeout"`
	Pinned         float64 `bson:"pinned"`
}

// Export exports the data to prometheus.
func (cursors *Cursors) Export() {
	cursorsMetric.WithLabelValues("total_open").Add(cursors.TotalOpen)
	cursorsMetric.WithLabelValues("timed_out").Add(cursors.TimeOut)
	cursorsMetric.WithLabelValues("total_no_timeout").Add(cursors.TotalNoTimeout)
	cursorsMetric.WithLabelValues("pinned").Add(cursors.Pinned)
}
