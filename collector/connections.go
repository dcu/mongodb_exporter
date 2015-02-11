package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"github.com/prometheus/client_golang/prometheus"
)

// Connection
type ConnectionStats struct {
	Current      float64 `bson:"current"`
	Available    float64 `bson:"available"`
	TotalCreated float64 `bson:"totalCreated"`
}

func (connectionStats *ConnectionStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)
	group.Collect("current", connectionStats.Current, ch)
	group.Collect("available", connectionStats.Available, ch)
	group.Collect("total_created", connectionStats.TotalCreated, ch)
}
