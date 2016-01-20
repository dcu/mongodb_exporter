package collector

import (
	//"github.com/dcu/mongodb_exporter/shared"
	"github.com/dcu/mongodb_exporter/collector/mongod"
	"github.com/dcu/mongodb_exporter/collector/mongos"
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

	return exporter
}

// Describe describes all mongodb's metrics.
func (exporter *MongodbCollector) Describe(ch chan<- *prometheus.Desc) {
	glog.Info("Describing groups")
	serverStatus := collector_mongod.GetServerStatus(exporter.Opts.URI)

	if serverStatus != nil {
		serverStatus.Describe(ch)
	}
}

// Collect collects all mongodb's metrics.
func (exporter *MongodbCollector) Collect(ch chan<- prometheus.Metric) {
	glog.Info("Collecting Server Status")
	/**
	We need to add logic here:
		if mongos:
		    collectMongosBalancingData
			collectMongosServerStatus
		else if replset
			collectMongodServerStatus
			collectElectionInfo
			collectOpLogInfo
			collectReplicationData
		else if mongod:
			collectMongodServerStatus
		else if arbiter:
			pass
		else:
			WTF()
	**/
	exporter.collectMongodServerStatus(ch)
	//exporter.collectMongosServerStatus(ch)
}

func (exporter *MongodbCollector) collectMongodServerStatus(ch chan<- prometheus.Metric) *collector_mongod.ServerStatus {
	serverStatus := collector_mongod.GetServerStatus(exporter.Opts.URI)

	if serverStatus != nil {
		serverStatus.Export(ch)
	}

	return serverStatus
}

func (exporter *MongodbCollector) collectMongosServerStatus(ch chan<- prometheus.Metric) *collector_mongos.ServerStatus {
	serverStatus := collector_mongos.GetServerStatus(exporter.Opts.URI)

	if serverStatus != nil {
		serverStatus.Export(ch)
	}

	return serverStatus
}

