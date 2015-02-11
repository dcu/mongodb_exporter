package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"github.com/prometheus/client_golang/prometheus"
)

type AssertsStats struct {
	Regular   float64 `bson:"regular"`
	Warning   float64 `bson:"warning"`
	Msg       float64 `bson:"msg"`
	User      float64 `bson:"user"`
	Rollovers float64 `bson:"rollovers"`
}

func (asserts *AssertsStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)

	group.Collect("regular", asserts.Regular, ch)
	group.Collect("warning", asserts.Warning, ch)
	group.Collect("msg", asserts.Msg, ch)
	group.Collect("user", asserts.User, ch)
	group.Collect("rollovers", asserts.Rollovers, ch)
}
