package collector

import(
    "github.com/prometheus/client_golang/prometheus"
    "github.com/dcu/mongodb_exporter/shared"
)

//Network
type NetworkStats struct {
    BytesIn             float64 `bson:"bytesIn" type:"gauge"`
    BytesOut            float64 `bson:"bytesOut" type:"gauge"`
    NumRequests         float64 `bson:"numRequests" type:"gauge"`
}

func (networkStats *NetworkStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)
    group.Collect(networkStats, "BytesIn", ch)
    group.Collect(networkStats, "BytesOut", ch)
    group.Collect(networkStats, "NumRequests", ch)
}

