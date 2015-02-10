package collector

import(
    "github.com/prometheus/client_golang/prometheus"
    "github.com/dcu/mongodb_exporter/shared"
)

type MongodbCollector struct {
}

func NewMongodbCollector() *MongodbCollector {
    exporter := &MongodbCollector{}
    exporter.collectServerStatus(nil)

    return exporter
}

func (exporter *MongodbCollector) Describe(ch chan<- *prometheus.Desc) {
    for _, group := range shared.Groups {
        group.Describe(ch)
    }
}

func (exporter *MongodbCollector) Collect(ch chan<- prometheus.Metric) {
    println("Collecting Server Status")
    exporter.collectServerStatus(ch)
}

func (exporter *MongodbCollector) collectServerStatus(ch chan<-prometheus.Metric) {
    serverStatus := GetServerStatus()
    serverStatus.Collect("instance", ch)
}

