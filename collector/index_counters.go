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
    group.Collect("accesses", connectionStats.Accesses, ch)
    group.Collect("hits", connectionStats.Hits, ch)
    group.Collect("misses", connectionStats.Misses, ch)
    group.Collect("resets", connectionStats.Resets, ch)
    group.Collect("miss_ratio", connectionStats.MissRatio, ch)
}
