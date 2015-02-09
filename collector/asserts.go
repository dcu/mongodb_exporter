package collector

import(
    "github.com/prometheus/client_golang/prometheus"
)

type AssertsStats struct {
    Regular float64 `bson:"regular" type:"counter"`
    Warning float64 `bson:"warning" type:"counter"`
    Msg float64 `bson:"msg" type:"counter"`
    User float64 `bson:"user" type:"counter"`
    Rollovers float64 `bson:"rollovers" type:"counter"`
}

func (asserts *AssertsStats) Collect(groupName string, exporter *MongodbCollector, ch chan<- prometheus.Metric) {
    group := exporter.FindOrCreateGroupByName(groupName)

    group.Collect(asserts, "Regular", ch)
    group.Collect(asserts, "Warning", ch)
    group.Collect(asserts, "Msg", ch)
    group.Collect(asserts, "User", ch)
    group.Collect(asserts, "Rollovers", ch)
}

