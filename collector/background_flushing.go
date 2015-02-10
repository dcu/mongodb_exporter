package collector

import(
    "github.com/prometheus/client_golang/prometheus"
    "github.com/dcu/mongodb_exporter/shared"
    "time"
)


// Flush
type FlushStats struct {
    Flushes float64 `bson:"flushes" type:"counter"`
    TotalMs float64 `bson:"total_ms" type:"counter"`
    AverageMs float64 `bson:"average_ms" type:"gauge"`
    LastMs float64 `bson:"last_ms" type:"gauge"`
    LastFinished time.Time `bson:"last_finished" type:"gauge"`
}

func (flushStats *FlushStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)
    group.Collect(flushStats, "Flushes", ch)
    group.Collect(flushStats, "TotalMs", ch)
    group.Collect(flushStats, "AverageMs", ch)
    group.Collect(flushStats, "LastMs", ch)
    group.Collect(flushStats, "LastFinished", ch)
}


