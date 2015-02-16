package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
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

func (opCounters *OpcountersStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName+"_total")
	group.Export("insert", opCounters.Insert)
	group.Export("query", opCounters.Query)
	group.Export("update", opCounters.Update)
	group.Export("delete", opCounters.Delete)
	group.Export("getmore", opCounters.GetMore)
	group.Export("command", opCounters.Command)
}
