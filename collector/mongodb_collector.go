package collector

import(
    "github.com/prometheus/client_golang/prometheus"
)

var(
    connections = prometheus.NewGauge(prometheus.GaugeOpts{
        Namespace: "mongodb",
        Subsystem: "connections",
        Name:      "current",
        Help:      "The value of current corresponds to the number of connections to the database server from clients.",
    })
)

type MongodbCollector struct {
}

func NewMongodbCollector() *MongodbCollector {
    return &MongodbCollector{}
}

func (e *MongodbCollector) Describe(ch chan<- *prometheus.Desc) {
    connections.Describe(ch)
}

func (e *MongodbCollector) Collect(ch chan<- prometheus.Metric) {
    println("Collecting Server Status")
    serverStatus := NewServerStatus()

    connections.Set(serverStatus.Connections.Current)
    connections.Collect(ch)
}


