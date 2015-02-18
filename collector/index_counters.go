package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
)

//IndexCounter
type IndexCounterStats struct {
	Accesses  float64 `bson:"accesses`
	Hits      float64 `bson:"hits"`
	Misses    float64 `bson:"misses"`
	Resets    float64 `bson:"resets"`
	MissRatio float64 `bson:"missRatio"`
}

func (indexCountersStats *IndexCounterStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName+"_total")
	group.Export("accesses", indexCountersStats.Accesses)
	group.Export("hits", indexCountersStats.Hits)
	group.Export("misses", indexCountersStats.Misses)
	group.Export("resets", indexCountersStats.Resets)

	group = shared.FindOrCreateGroup(groupName)
	group.Export("miss_ratio", indexCountersStats.MissRatio)
}

