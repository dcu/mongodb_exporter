package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
)

// Connection
type ConnectionStats struct {
	Current      float64 `bson:"current"`
	Available    float64 `bson:"available"`
	TotalCreated float64 `bson:"totalCreated"`
}

func (connectionStats *ConnectionStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("current", connectionStats.Current)
	group.Export("available", connectionStats.Available)

	group = shared.FindOrCreateGroup(groupName+"_metrics")
	group.Export("created_total", connectionStats.TotalCreated)
}
