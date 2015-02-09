package collector

import(
    "github.com/prometheus/client_golang/prometheus"
)

type MongodbCollector struct {
    Groups map[string]*Group
}

func NewMongodbCollector() *MongodbCollector {
    exporter := &MongodbCollector{
        Groups: make(map[string]*Group),
    }

    exporter.collectServerStatus(nil)

    return exporter
}

func (exporter *MongodbCollector) FindOrCreateGroupByName(name string) *Group {
    name = SnakeCase(name)
    println("Adding group:",name)
    group := exporter.Groups[name]

    if group == nil {
        group = NewGroup(name)
        exporter.Groups[name] = group
    }

    return group
}

func (exporter *MongodbCollector) Describe(ch chan<- *prometheus.Desc) {
    for _, group := range exporter.Groups {
        group.Describe(ch)
    }
}

func (exporter *MongodbCollector) Collect(ch chan<- prometheus.Metric) {
    println("Collecting Server Status")
    exporter.collectServerStatus(ch)
}

func (exporter *MongodbCollector) collectServerStatus(ch chan<-prometheus.Metric) {
    serverStatus := GetServerStatus()
    serverStatus.Collect("instance", exporter, ch)
}

