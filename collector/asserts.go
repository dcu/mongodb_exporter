package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
)

type AssertsStats struct {
	Regular   float64 `bson:"regular"`
	Warning   float64 `bson:"warning"`
	Msg       float64 `bson:"msg"`
	User      float64 `bson:"user"`
	Rollovers float64 `bson:"rollovers"`
}

func (asserts *AssertsStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName+"_total")

	group.Export("regular", asserts.Regular)
	group.Export("warning", asserts.Warning)
	group.Export("msg", asserts.Msg)
	group.Export("user", asserts.User)
	group.Export("rollovers", asserts.Rollovers)
}

