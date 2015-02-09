package main

import(
    "net/http"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/dcu/mongodb_exporter/collector"
)

func main() {
    mongodbCollector := collector.NewMongodbCollector()
    prometheus.MustRegister(mongodbCollector)

    http.Handle("/metrics", prometheus.Handler())
    http.ListenAndServe(":9001", nil)
}

