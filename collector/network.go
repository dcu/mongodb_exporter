package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
)

//Network
type NetworkStats struct {
	BytesIn     float64 `bson:"bytesIn"`
	BytesOut    float64 `bson:"bytesOut"`
	NumRequests float64 `bson:"numRequests"`
}

func (networkStats *NetworkStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName + "_bytes_total")
	group.Export("in_bytes", networkStats.BytesIn)
	group.Export("out_bytes", networkStats.BytesOut)

	group = shared.FindOrCreateGroup(groupName + "_metrics")
	group.Export("num_requests_total", networkStats.NumRequests)
}
