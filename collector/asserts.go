package collector

import(
    "github.com/prometheus/client_golang/prometheus"
    "github.com/dcu/mongodb_exporter/shared"
)

type AssertsStats struct {
    Regular float64 `bson:"regular" type:"counter"`
    Warning float64 `bson:"warning" type:"counter"`
    Msg float64 `bson:"msg" type:"counter"`
    User float64 `bson:"user" type:"counter"`
    Rollovers float64 `bson:"rollovers" type:"counter"`
}

func (asserts *AssertsStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)

    group.Collect("regular", asserts.Regular, ch)
    group.Collect("warning", asserts.Warning, ch)
    group.Collect("msg", asserts.Msg, ch)
    group.Collect("user", asserts.User, ch)
    group.Collect("rollovers", asserts.Rollovers, ch)
}

