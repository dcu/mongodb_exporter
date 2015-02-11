package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"github.com/prometheus/client_golang/prometheus"
)

//Opcount and OpcountersRepl
type OpcountersStats struct {
	Insert  float64 `bson:"insert"`
	Query   float64 `bson:"query"`
	Update  float64 `bson:"update"`
	Delete  float64 `bson:"delete"`
	GetMore float64 `bson:"getmore"`
	Command float64 `bson:"command"`
}

func (opCounters *OpcountersStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)
	group.Collect("insert", opCounters.Insert, ch)
	group.Collect("query", opCounters.Query, ch)
	group.Collect("update", opCounters.Update, ch)
	group.Collect("delete", opCounters.Delete, ch)
	group.Collect("getmore", opCounters.GetMore, ch)
	group.Collect("command", opCounters.Command, ch)
}
