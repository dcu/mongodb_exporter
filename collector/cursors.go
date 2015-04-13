package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
)

type Cursors struct {
	TotalOpen      float64 `bson:"totalOpen"`
	TimeOut        float64 `bson:"timedOut"`
	TotalNoTimeout float64 `bson:"totalNoTimeout"`
	Pinned         float64 `bson:"pinned"`
}

func (cursors *Cursors) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)

	group.Export("total_open", cursors.TotalOpen)
	group.Export("timed_out", cursors.TimeOut)
	group.Export("total_no_timeout", cursors.TotalNoTimeout)
	group.Export("pinned", cursors.Pinned)
}
