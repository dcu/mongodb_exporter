package collector

import (
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

	Opcounters     *OpcountersStats     `bson:"opcounters"`
	OpcountersRepl *OpcountersReplStats `bson:"opcountersRepl"`
	Mem            *MemStats            `bson:"mem"`
	Metrics        *MetricsStats        `bson:"metrics"`

	Cursors *Cursors `bson:"cursors"`
}

// Export exports the server status to be consumed by prometheus.
func (status *ServerStatus) Export() {

	instanceUptimeSeconds.Set(status.Uptime)
	instanceUptimeEstimateSeconds.Set(status.Uptime)
	instanceLocalTime.Set(float64(status.LocalTime.Unix()))

	if status.Asserts != nil {
		status.Asserts.Export()
	}
	if status.Dur != nil {
		status.Dur.Export()
	}
	if status.BackgroundFlushing != nil {
		status.BackgroundFlushing.Export()
	}
	if status.Connections != nil {
		status.Connections.Export()
	}
	if status.Asserts != nil {
		status.Asserts.Export()
	}
	if status.ExtraInfo != nil {
		status.ExtraInfo.Export()
	}
	if status.GlobalLock != nil {
		status.GlobalLock.Export()
	}
	if status.IndexCounter != nil {
		status.IndexCounter.Export()
	}
	if status.Network != nil {
		status.Network.Export()
	}
	if status.Opcounters != nil {
		status.Opcounters.Export()
	}
	if status.OpcountersRepl != nil {
		status.OpcountersRepl.Export()
	}
	if status.Mem != nil {
		status.Mem.Export()
	}
	if status.Locks != nil {
		status.Locks.Export()
	}
	if status.Metrics != nil {
		status.Metrics.Export()
	}
	if status.Cursors != nil {
		status.Cursors.Export()
	}
}

// Describe describes the server status for prometheus.
func (status *ServerStatus) Describe(ch chan<- *prometheus.Desc) {
	instanceUptimeSeconds.Describe(ch)
	instanceUptimeEstimateSeconds.Describe(ch)
	instanceLocalTime.Describe(ch)

	if status.Asserts != nil {
		status.Asserts.Describe(ch)
	}
	if status.Dur != nil {
		status.Dur.Describe(ch)
	}
	if status.BackgroundFlushing != nil {
		status.BackgroundFlushing.Describe(ch)
	}
	if status.Connections != nil {
		status.Connections.Describe(ch)
	}
	if status.ExtraInfo != nil {
		status.ExtraInfo.Describe(ch)
	}
	if status.GlobalLock != nil {
		status.GlobalLock.Describe(ch)
	}
	if status.IndexCounter != nil {
		status.IndexCounter.Describe(ch)
	}
	if status.Network != nil {
		status.Network.Describe(ch)
	}
	if status.Opcounters != nil {
		status.Opcounters.Describe(ch)
	}
	if status.OpcountersRepl != nil {
		status.OpcountersRepl.Describe(ch)
	}
	if status.Mem != nil {
		status.Mem.Describe(ch)
	}
	if status.Locks != nil {
		status.Locks.Describe(ch)
	}
	if status.Metrics != nil {
		status.Metrics.Describe(ch)
	}
	if status.Cursors != nil {
		status.Cursors.Describe(ch)
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
