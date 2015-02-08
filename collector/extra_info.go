package collector

import(
    "github.com/prometheus/client_golang/prometheus"
)


// ExtraInfo
type ExtraInfo struct {
    HeapUsageBytes float64 `bson:"heap_usage_bytes" type:"gauge"`
    PageFaults float64 `bson:"page_faults" type:"gauge"`
}

func (extraInfo *ExtraInfo) Collect(exporter *MongodbCollector, ch chan<- prometheus.Metric) {
    group := exporter.FindOrCreateGroupByName("extra_info")
    group.Collect(extraInfo, "HeapUsageBytes", ch)
    group.Collect(extraInfo, "PageFaults", ch)
}

