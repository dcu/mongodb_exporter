package collector

import(
    "github.com/prometheus/client_golang/prometheus"
)


// GlobalLock
type ClientStats struct {
    Total float64 `bson:"total" type:"gauge"`
    Readers float64 `bson:"readers" type:"gauge"`
    Writers float64 `bson:"writers" type:"gauge"`
}
func (clientStats *ClientStats) Collect(groupName string, exporter *MongodbCollector, ch chan<- prometheus.Metric) {
    group := exporter.FindOrCreateGroupByName(groupName)
    group.Collect(clientStats, "Total", ch)
    group.Collect(clientStats, "Readers", ch)
    group.Collect(clientStats, "Writers", ch)
}

type QueueStats struct {
    Total float64 `bson:"total" type:"gauge"`
    Readers float64 `bson:"readers" type:"gauge"`
    Writers float64 `bson:"writers" type:"gauge"`
}
func (queueStats *QueueStats) Collect(groupName string, exporter *MongodbCollector, ch chan<- prometheus.Metric) {
    group := exporter.FindOrCreateGroupByName(groupName)
    group.Collect(queueStats, "Total", ch)
    group.Collect(queueStats, "Readers", ch)
    group.Collect(queueStats, "Writers", ch)
}

type GlobalLockStats struct {
    TotalTime float64 `bson:"totalTime" type:"counter"`
    LockTime float64 `bson:"lockTime" type:"counter"`
    Ratio float64 `bson:"ratio" type:"gauge"`
    CurrentQueue *QueueStats `bson:"currentQueue"`
    ActiveClients *ClientStats `bson:"activeClients"`
}
func (globalLock *GlobalLockStats) Collect(groupName string, exporter *MongodbCollector, ch chan<- prometheus.Metric) {
    group := exporter.FindOrCreateGroupByName(groupName)
    group.Collect(globalLock, "TotalTime", ch)
    group.Collect(globalLock, "LockTime", ch)
    group.Collect(globalLock, "Ratio", ch)

    globalLock.CurrentQueue.Collect(groupName+"_queue", exporter, ch)
    globalLock.ActiveClients.Collect(groupName+"_client", exporter, ch)
}


