package collector

import(
    "github.com/prometheus/client_golang/prometheus"
)

//Network
type NetworkStats struct {
    BytesIn             float64 `bson:"bytesIn" type:"gauge"`
    BytesOut            float64 `bson:"bytesOut" type:"gauge"`
    NumRequests         float64 `bson:"numRequests" type:"gauge"`
}

func (networkStats *NetworkStats) Collect(exporter *MongodbCollector, ch chan<- prometheus.Metric) {
    group := exporter.FindOrCreateGroupByName("network")
    group.Collect(networkStats, "BytesIn", ch)
    group.Collect(networkStats, "BytesOut", ch)
    group.Collect(networkStats, "NumRequests", ch)
}

