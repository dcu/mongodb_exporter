package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/golang/glog"
)

type MongodbCollectorOpts struct {
	URI string
}

type MongodbCollector struct {
	Opts MongodbCollectorOpts
}

func NewMongodbCollector(opts MongodbCollectorOpts) *MongodbCollector {
	exporter := &MongodbCollector{
		Opts: opts,
	}
	exporter.collectServerStatus(nil)

	return exporter
}

func (exporter *MongodbCollector) Describe(ch chan<- *prometheus.Desc) {
	glog.Info("Describing groups")
	for _, group := range shared.Groups {
		group.Describe(ch)
	}
}

func (exporter *MongodbCollector) Collect(ch chan<- prometheus.Metric) {
	glog.Info("Collecting Server Status")
	exporter.collectServerStatus(ch)
}

func (exporter *MongodbCollector) collectServerStatus(ch chan<- prometheus.Metric) *ServerStatus {
	serverStatus := GetServerStatus(exporter.Opts.URI)
	serverStatus.Collect("instance", ch)

	return serverStatus
}
