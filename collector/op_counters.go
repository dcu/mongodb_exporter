package collector

import(
    "github.com/prometheus/client_golang/prometheus"
)

//Opcount and OpcountersRepl
type OpcountersStats struct {
    Insert  float64 `bson:"insert" type:"gauge"`
    Query   float64 `bson:"query" type:"gauge"`
    Update  float64 `bson:"update" type:"gauge"`
    Delete  float64 `bson:"delete" type:"gauge"`
    GetMore float64 `bson:"getmore" type:"gauge"`
    Command float64 `bson:"command" type:"gauge"`
}

func (opCounters *OpcountersStats) Collect(groupName string, exporter *MongodbCollector, ch chan<- prometheus.Metric) {
    group := exporter.FindOrCreateGroupByName(groupName)
    group.Collect(opCounters, "Insert", ch)
    group.Collect(opCounters, "Query", ch)
    group.Collect(opCounters, "Update", ch)
    group.Collect(opCounters, "Delete", ch)
    group.Collect(opCounters, "GetMore", ch)
    group.Collect(opCounters, "Command", ch)
}

