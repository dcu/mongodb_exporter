package collector

import (
	//"github.com/dcu/mongodb_exporter/shared"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Namespace is the namespace of the metrics
	Namespace = "mongodb"
)

// MongodbCollectorOpts is the options of the mongodb collector.
type MongodbCollectorOpts struct {
	URI string
}

// MongodbCollector is in charge of collecting mongodb's metrics.
type MongodbCollector struct {
	Opts MongodbCollectorOpts
}

// NewMongodbCollector returns a new instance of a MongodbCollector.
func NewMongodbCollector(opts MongodbCollectorOpts) *MongodbCollector {
	exporter := &MongodbCollector{
		Opts: opts,
	}
	exporter.collectServerStatus(nil)

	return exporter
}

// Describe describes all mongodb's metrics.
func (exporter *MongodbCollector) Describe(ch chan<- *prometheus.Desc) {
	glog.Info("Describing groups")
	serverStatus := GetServerStatus(exporter.Opts.URI)

	if serverStatus != nil {
		serverStatus.Describe(ch)
	}
}

// Collect collects all mongodb's metrics.
func (exporter *MongodbCollector) Collect(ch chan<- prometheus.Metric) {
	glog.Info("Collecting Server Status")
	exporter.collectServerStatus(ch)
}

func (exporter *MongodbCollector) collectServerStatus(ch chan<- prometheus.Metric) *ServerStatus {
	serverStatus := GetServerStatus(exporter.Opts.URI)

	if serverStatus != nil {
		serverStatus.Export()
	}

	return serverStatus
}
