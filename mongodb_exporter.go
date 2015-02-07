package main

import(
    "net/http"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/dcu/mongodb_exporter/collector"
)

func main() {
    mongodb_collector := collector.NewMongodbCollector()
    mongodb_collector.LoadGroups()

    prometheus.MustRegister(mongodb_collector)

    collector.NewServerStatus()

    http.Handle("/metrics", prometheus.Handler())
    http.ListenAndServe(":9001", nil)
}

