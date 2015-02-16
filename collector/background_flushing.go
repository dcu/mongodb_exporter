package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"time"
)

// Flush
type FlushStats struct {
	Flushes      float64   `bson:"flushes"`
	TotalMs      float64   `bson:"total_ms"`
	AverageMs    float64   `bson:"average_ms"`
	LastMs       float64   `bson:"last_ms"`
	LastFinished time.Time `bson:"last_finished"`
}

func (flushStats *FlushStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("flushes_total", flushStats.Flushes)
	group.Export("total_milliseconds", flushStats.TotalMs)
	group.Export("average_milliseconds", flushStats.AverageMs)
	group.Export("last_milliseconds", flushStats.LastMs)
	group.Export("last_finished_time", float64(flushStats.LastFinished.Unix()))
}
