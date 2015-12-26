package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	instanceUptimeSeconds = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "instance",
		Name:      "uptime_seconds",
		Help:      "The value of the uptime field corresponds to the number of seconds that the mongos or mongod process has been active.",
	})
	instanceUptimeEstimateSeconds = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "instance",
		Name:      "uptime_estimate_seconds",
		Help:      "uptimeEstimate provides the uptime as calculated from MongoDB's internal course-grained time keeping system.",
	})
	instanceLocalTime = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "instance",
		Name:      "local_time",
		Help:      "The localTime value is the current time, according to the server, in UTC specified in an ISODate format.",
	})
)

// ServerStatus keeps the data returned by the serverStatus() method.
type ServerStatus struct {
	Uptime         float64   `bson:"uptime"`
	UptimeEstimate float64   `bson:"uptimeEstimate"`
	LocalTime      time.Time `bson:"localTime"`

	Asserts *AssertsStats `bson:"asserts"`

	Dur *DurStats `bson:"dur"`

	BackgroundFlushing *FlushStats `bson:"backgroundFlushing"`

	Connections *ConnectionStats `bson:"connections"`

	ExtraInfo *ExtraInfo `bson:"extra_info"`

	GlobalLock *GlobalLockStats `bson:"globalLock"`

	IndexCounter *IndexCounterStats `bson:"indexCounters"`

	Locks LockStatsMap `bson:"locks,omitempty"`

	Network *NetworkStats `bson:"network"`

	Opcounters     *OpcountersStats `bson:"opcounters"`
	OpcountersRepl *OpcountersStats `bson:"opcountersRepl"`
	Mem            *MemStats        `bson:"mem"`
	Metrics        *MetricsStats    `bson:"metrics"`

	Cursors *Cursors `bson:"cursors"`
}

// Export exports the given groupName to be consumed by prometheus.
func (status *ServerStatus) Export() {

	instanceUptimeSeconds.Set(status.Uptime)
	instanceUptimeEstimateSeconds.Set(status.Uptime)
	instanceLocalTime.Set(float64(status.LocalTime.Unix()))

	exportData(status.Asserts)
	exportData(status.Dur)
	exportData(status.BackgroundFlushing)
	exportData(status.Connections)
	exportData(status.ExtraInfo)
	exportData(status.GlobalLock)
	exportData(status.IndexCounter)
	exportData(status.Network)
	exportData(status.Opcounters)
	exportData(status.OpcountersRepl)
	exportData(status.Mem)
	exportData(status.Locks)
	exportData(status.Metrics)
	exportData(status.Cursors)
}

func exportData(exportable shared.Exportable) {
	if exportable != nil {
		exportable.Export()
	}
}

// GetServerStatus returns the server status info.
func GetServerStatus(uri string) *ServerStatus {
	result := &ServerStatus{}
	session, err := mgo.Dial(uri)
	if err != nil {
		glog.Errorf("Cannot connect to server using url: %s", uri)
		return nil
	}

	session.SetMode(mgo.Eventual, true)
	session.SetSocketTimeout(0)
	defer func() {
		glog.Info("Closing connection to database.")
		session.Close()
	}()

	err = session.DB("admin").Run(bson.D{{"serverStatus", 1}, {"recordStats", 0}}, result)
	if err != nil {
		glog.Error("Failed to get server status.")
		return nil
	}

	return result
}
