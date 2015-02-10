package collector

import(
    "github.com/prometheus/client_golang/prometheus"
    "github.com/dcu/mongodb_exporter/shared"
)

//Mem
type MemStats struct {
    Bits                float64 `bson:"bits" type:"counter"`
    Resident            float64 `bson:"resident" type:"counter"`
    Virtual             float64 `bson:"virtual" type:"gauge"`
    Mapped              float64 `bson:"mapped" type:"gauge"`
    MappedWithJournal   float64 `bson:"mappedWithJournal" type:"counter"`
}

func (memStats *MemStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)
    group.Collect("bits", memStats.Bits, ch)
    group.Collect("resident", memStats.Resident, ch)
    group.Collect("virtual", memStats.Virtual, ch)
    group.Collect("mapped", memStats.Mapped, ch)
    group.Collect("mapped_with_journal", memStats.MappedWithJournal, ch)
}


