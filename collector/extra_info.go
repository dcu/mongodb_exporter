package collector

import(
    "github.com/prometheus/client_golang/prometheus"
    "github.com/dcu/mongodb_exporter/shared"
)


// ExtraInfo
type ExtraInfo struct {
    HeapUsageBytes float64 `bson:"heap_usage_bytes" type:"gauge"`
    PageFaults float64 `bson:"page_faults" type:"gauge"`
}

func (extraInfo *ExtraInfo) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)
    group.Collect("heap_usage_bytes", extraInfo.HeapUsageBytes, ch)
    group.Collect("page_faults", extraInfo.PageFaults, ch)
}

