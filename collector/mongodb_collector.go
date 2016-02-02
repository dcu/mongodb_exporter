package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/mgo.v2"
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
	(&ServerStatus{}).Describe(ch)
	(&ReplSetStatus{}).Describe(ch)
}

// Collect collects all mongodb's metrics.
func (exporter *MongodbCollector) Collect(ch chan<- prometheus.Metric) {
	mongoSess := shared.MongoSession(exporter.Opts.URI)
	defer mongoSess.Close()
	if mongoSess != nil {
		glog.Info("Collecting Server Status")
		exporter.collectServerStatus(mongoSess, ch)
		glog.Info("Collecting ReplSet Status")
		exporter.collectReplSetStatus(mongoSess, ch)
	}
}

func (exporter *MongodbCollector) collectServerStatus(session *mgo.Session, ch chan<- prometheus.Metric) *ServerStatus {
	serverStatus := GetServerStatus(session)
	if serverStatus != nil {
		glog.Info("exporting ServerStatus Metrics")
		serverStatus.Export(ch)
	}
	return serverStatus
}

func (exporter *MongodbCollector) collectReplSetStatus(session *mgo.Session, ch chan<- prometheus.Metric) *ReplSetStatus {
	replSetStatus := GetReplSetStatus(session)

	if replSetStatus != nil {
		glog.Info("exporting ReplSetStatus Metrics")
		replSetStatus.Export(ch)
	}

	return replSetStatus
}
