package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"github.com/prometheus/client_golang/prometheus"
)

//Mem
type MemStats struct {
	Bits              float64 `bson:"bits"`
	Resident          float64 `bson:"resident"`
	Virtual           float64 `bson:"virtual"`
	Mapped            float64 `bson:"mapped"`
	MappedWithJournal float64 `bson:"mappedWithJournal"`
}

func (memStats *MemStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)
	group.Collect("bits", memStats.Bits, ch)
	group.Collect("resident", memStats.Resident, ch)
	group.Collect("virtual", memStats.Virtual, ch)
	group.Collect("mapped", memStats.Mapped, ch)
	group.Collect("mapped_with_journal", memStats.MappedWithJournal, ch)
}
