package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
)

// ExtraInfo
type ExtraInfo struct {
	HeapUsageBytes float64 `bson:"heap_usage_bytes"`
	PageFaults     float64 `bson:"page_faults"`
}

func (extraInfo *ExtraInfo) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("heap_usage_bytes", extraInfo.HeapUsageBytes)
	group.Export("page_faults_total", extraInfo.PageFaults)
}
