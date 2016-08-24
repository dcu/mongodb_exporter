package collector

import (
	"os"

	"github.com/dcu/mongodb_exporter/shared"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/mgo.v2"
)

var (
	// Namespace is the namespace of the metrics
	Namespace = namespace()
)

func namespace() string {
	n := os.Getenv("NAMESPACE")
	if n == "" {
		return "mongodb"
	}
	return n
}

// MongodbCollectorOpts is the options of the mongodb collector.
type MongodbCollectorOpts struct {
	URI                   string
	TLSCertificateFile    string
	TLSPrivateKeyFile     string
	TLSCaFile             string
	TLSHostnameValidation bool
	CollectReplSet        bool
	CollectOplog          bool
}

func (in MongodbCollectorOpts) toSessionOps() shared.MongoSessionOpts {
	return shared.MongoSessionOpts{
		URI:                   in.URI,
		TLSCertificateFile:    in.TLSCertificateFile,
		TLSPrivateKeyFile:     in.TLSPrivateKeyFile,
		TLSCaFile:             in.TLSCaFile,
		TLSHostnameValidation: in.TLSHostnameValidation,
	}
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
	mongoSess := shared.MongoSession(exporter.Opts.toSessionOps())
	if mongoSess != nil {
		defer mongoSess.Close()
		glog.Info("Collecting Server Status")
		exporter.collectServerStatus(mongoSess, ch)
		if exporter.Opts.CollectReplSet {
			glog.Info("Collecting ReplSet Status")
			exporter.collectReplSetStatus(mongoSess, ch)
		}
		if exporter.Opts.CollectOplog {
			glog.Info("Collecting Oplog Status")
			exporter.collectOplogStatus(mongoSess, ch)
		}
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

func (exporter *MongodbCollector) collectOplogStatus(session *mgo.Session, ch chan<- prometheus.Metric) *OplogStatus {
	oplogStatus := GetOplogStatus(session)

	if oplogStatus != nil {
		glog.Info("exporting OplogStatus Metrics")
		oplogStatus.Export(ch)
	}

	return oplogStatus
}
