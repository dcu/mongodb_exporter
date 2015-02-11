package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"github.com/prometheus/client_golang/prometheus"
)

// ExtraInfo
type ExtraInfo struct {
	HeapUsageBytes float64 `bson:"heap_usage_bytes"`
	PageFaults     float64 `bson:"page_faults"`
}

func (extraInfo *ExtraInfo) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)
	group.Collect("heap_usage_bytes", extraInfo.HeapUsageBytes, ch)
	group.Collect("page_faults", extraInfo.PageFaults, ch)
}
