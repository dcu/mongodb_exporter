package main

import(
    "net/http"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/dcu/mongodb_exporter/collector"
    "github.com/dcu/mongodb_exporter/shared"
)

func main() {
    shared.LoadGroupsDesc()

    mongodbCollector := collector.NewMongodbCollector()
    prometheus.MustRegister(mongodbCollector)

    http.Handle("/metrics", prometheus.Handler())
    err := http.ListenAndServe(":9001", nil)

    if err != nil {
        panic(err)
    }
}

