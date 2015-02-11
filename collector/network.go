package collector

import(
    "github.com/prometheus/client_golang/prometheus"
    "github.com/dcu/mongodb_exporter/shared"
)

//Network
type NetworkStats struct {
    BytesIn             float64 `bson:"bytesIn"`
    BytesOut            float64 `bson:"bytesOut"`
    NumRequests         float64 `bson:"numRequests"`
}

func (networkStats *NetworkStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)
    group.Collect("bytes_in", networkStats.BytesIn, ch)
    group.Collect("bytes_out", networkStats.BytesOut, ch)
    group.Collect("num_requests", networkStats.NumRequests, ch)
}

