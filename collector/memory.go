package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
)

//Mem
type MemStats struct {
	Bits              float64 `bson:"bits"`
	Resident          float64 `bson:"resident"`
	Virtual           float64 `bson:"virtual"`
	Mapped            float64 `bson:"mapped"`
	MappedWithJournal float64 `bson:"mappedWithJournal"`
}

func (memStats *MemStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("resident", memStats.Resident)
	group.Export("virtual", memStats.Virtual)
	group.Export("mapped", memStats.Mapped)
	group.Export("mapped_with_journal", memStats.MappedWithJournal)
}
