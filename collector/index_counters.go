package collector

import(
    "github.com/prometheus/client_golang/prometheus"
    "github.com/dcu/mongodb_exporter/shared"
)


//IndexCounter
type IndexCounterStats struct {
    Accesses float64 `bson:"accesses type:"counter"`
    Hits float64 `bson:"hits" type:"counter"`
    Misses float64 `bson:"misses" type:"counter"`
    Resets float64 `bson:"resets" type:"gauge"`
    MissRatio float64 `bson:"missRatio" type:"gauge"`
}
func (connectionStats *IndexCounterStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)
    group.Collect(connectionStats, "Accesses", ch)
    group.Collect(connectionStats, "Hits", ch)
    group.Collect(connectionStats, "Misses", ch)
    group.Collect(connectionStats, "Resets", ch)
    group.Collect(connectionStats, "MissRatio", ch)
}
