package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
)

// GlobalLock
type ClientStats struct {
	Total   float64 `bson:"total"`
	Readers float64 `bson:"readers"`
	Writers float64 `bson:"writers"`
}

func (clientStats *ClientStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("reader", clientStats.Readers)
	group.Export("writer", clientStats.Writers)
}

type QueueStats struct {
	Total   float64 `bson:"total"`
	Readers float64 `bson:"readers"`
	Writers float64 `bson:"writers"`
}

func (queueStats *QueueStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("reader", queueStats.Readers)
	group.Export("writer", queueStats.Writers)
}

type GlobalLockStats struct {
	TotalTime     float64      `bson:"totalTime"`
	LockTime      float64      `bson:"lockTime"`
	Ratio         float64      `bson:"ratio"`
	CurrentQueue  *QueueStats  `bson:"currentQueue"`
	ActiveClients *ClientStats `bson:"activeClients"`
}

func (globalLock *GlobalLockStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("lock_total", globalLock.LockTime)
	group.Export("ratio", globalLock.Ratio)

	globalLock.CurrentQueue.Export(groupName + "_current_queue")
	globalLock.ActiveClients.Export(groupName + "_client")
}
