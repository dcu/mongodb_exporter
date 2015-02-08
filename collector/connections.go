package collector

import(
    "github.com/prometheus/client_golang/prometheus"
)


// Connection
type ConnectionStats struct {
    Current float64 `bson:"current" type:"gauge"`
    Available float64 `bson:"available" type:"gauge"`
    TotalCreated float64 `bson:"totalCreated" type:"counter"`
}

func (connectionStats *ConnectionStats) Collect(exporter *MongodbCollector, ch chan<- prometheus.Metric) {
    group := exporter.FindOrCreateGroupByName("connections")
    group.Collect(connectionStats, "Current", ch)
    group.Collect(connectionStats, "Available", ch)
    group.Collect(connectionStats, "TotalCreated", ch)
}


