package collector

import(
	"github.com/dcu/mongodb_exporter/shared"
	"github.com/prometheus/client_golang/prometheus"
)

type Cursors struct {
	TotalOpen float64 `bson:"totalOpen"`
	TimeOut float64 `bson:"timedOut"`
	TotalNoTimeout float64 `bson:"totalNoTimeout"`
	Pinned float64 `bson:"pinned"`
}

func (cursors *Cursors) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)

	group.Collect("total_open", cursors.TotalOpen, ch)
	group.Collect("timed_out", cursors.TimeOut, ch)
	group.Collect("total_no_timeout", cursors.TotalNoTimeout, ch)
	group.Collect("pinned", cursors.Pinned, ch)
}

