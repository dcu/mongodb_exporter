package main

import(
    "net/http"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/dcu/mongodb_exporter/collector"
)

func main() {
    a := collector.NewMongodbCollector()
    prometheus.MustRegister(a)

    http.Handle("/metrics", prometheus.Handler())
    http.ListenAndServe(":9001", nil)
}
