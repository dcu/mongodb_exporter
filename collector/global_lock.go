package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"github.com/prometheus/client_golang/prometheus"
)

// GlobalLock
type ClientStats struct {
	Total   float64 `bson:"total"`
	Readers float64 `bson:"readers"`
	Writers float64 `bson:"writers"`
}

func (clientStats *ClientStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)
	group.Collect("total", clientStats.Total, ch)
	group.Collect("readers", clientStats.Readers, ch)
	group.Collect("writers", clientStats.Writers, ch)
}

type QueueStats struct {
	Total   float64 `bson:"total"`
	Readers float64 `bson:"readers"`
	Writers float64 `bson:"writers"`
}

func (queueStats *QueueStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)
	group.Collect("total", queueStats.Total, ch)
	group.Collect("readers", queueStats.Readers, ch)
	group.Collect("writers", queueStats.Writers, ch)
}

type GlobalLockStats struct {
	TotalTime     float64      `bson:"totalTime"`
	LockTime      float64      `bson:"lockTime"`
	Ratio         float64      `bson:"ratio"`
	CurrentQueue  *QueueStats  `bson:"currentQueue"`
	ActiveClients *ClientStats `bson:"activeClients"`
}

func (globalLock *GlobalLockStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)
	group.Collect("total_time", globalLock.TotalTime, ch)
	group.Collect("lock_time", globalLock.LockTime, ch)
	group.Collect("ratio", globalLock.Ratio, ch)

	globalLock.CurrentQueue.Collect(groupName+"_current_queue", ch)
	globalLock.ActiveClients.Collect(groupName+"_client", ch)
}
