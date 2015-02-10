package collector

import(
    "github.com/prometheus/client_golang/prometheus"
    "github.com/dcu/mongodb_exporter/shared"
)


// Connection
type ConnectionStats struct {
    Current float64 `bson:"current" type:"gauge"`
    Available float64 `bson:"available" type:"gauge"`
    TotalCreated float64 `bson:"totalCreated" type:"counter"`
}

func (connectionStats *ConnectionStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)
    group.Collect(connectionStats, "Current", ch)
    group.Collect(connectionStats, "Available", ch)
    group.Collect(connectionStats, "TotalCreated", ch)
}


