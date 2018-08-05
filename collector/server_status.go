package collector

import (
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	instanceUptimeSeconds = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "instance",
		Name:      "uptime_seconds",
		Help:      "The value of the uptime field corresponds to the number of seconds that the mongos or mongod process has been active.",
	})
	instanceUptimeEstimateSeconds = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "instance",
		Name:      "uptime_estimate_seconds",
		Help:      "uptimeEstimate provides the uptime as calculated from MongoDB's internal course-grained time keeping system.",
	})
	instanceLocalTime = prometheus.NewGauge(prometheus.GaugeOpts{
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
	TCMallocStats  *TCMallocStats       `bson:"tcmalloc"`
	Mem            *MemStats            `bson:"mem"`
	Metrics        *MetricsStats        `bson:"metrics"`

	Cursors *Cursors `bson:"cursors"`

	StorageEngine *StorageEngineStats `bson:"storageEngine"`
	WiredTiger    *WiredTigerStats    `bson:"wiredTiger"`
}

// Export exports the server status to be consumed by prometheus.
func (status *ServerStatus) Export(ch chan<- prometheus.Metric) {
	instanceUptimeSeconds.Set(status.Uptime)
	instanceUptimeEstimateSeconds.Set(status.Uptime)
	instanceLocalTime.Set(float64(status.LocalTime.Unix()))
	instanceUptimeSeconds.Collect(ch)
	instanceUptimeEstimateSeconds.Collect(ch)
	instanceLocalTime.Collect(ch)

	if status.Asserts != nil {
		status.Asserts.Export(ch)
	}
	if status.Dur != nil {
		status.Dur.Export(ch)
	}
	if status.BackgroundFlushing != nil {
		status.BackgroundFlushing.Export(ch)
	}
	if status.Connections != nil {
		status.Connections.Export(ch)
	}
	if status.ExtraInfo != nil {
		status.ExtraInfo.Export(ch)
	}
	if status.GlobalLock != nil {
		status.GlobalLock.Export(ch)
	}
	if status.IndexCounter != nil {
		status.IndexCounter.Export(ch)
	}
	if status.Network != nil {
		status.Network.Export(ch)
	}
	if status.Opcounters != nil {
		status.Opcounters.Export(ch)
	}
	if status.OpcountersRepl != nil {
		status.OpcountersRepl.Export(ch)
	}
	if status.TCMallocStats != nil {
		status.TCMallocStats.Export(ch)
	}
	if status.Mem != nil {
		status.Mem.Export(ch)
	}
	if status.Locks != nil {
		status.Locks.Export(ch)
	}
	if status.Metrics != nil {
		status.Metrics.Export(ch)
	}
	if status.Cursors != nil {
		status.Cursors.Export(ch)
	}

	if status.WiredTiger != nil {
		status.WiredTiger.Export(ch)
	}
	// If db.serverStatus().storageEngine does not exist (3.0+ only) and status.BackgroundFlushing does (MMAPv1 only), default to mmapv1
	// https://docs.mongodb.com/v3.0/reference/command/serverStatus/#storageengine
	if status.StorageEngine == nil && status.BackgroundFlushing != nil {
		status.StorageEngine = &StorageEngineStats{
			Name: "mmapv1",
		}
	}
	if status.StorageEngine != nil {
		status.StorageEngine.Export(ch)
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
	if status.TCMallocStats != nil {
		status.TCMallocStats.Describe(ch)
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
	if status.WiredTiger != nil {
		status.WiredTiger.Describe(ch)
	}
	if status.StorageEngine != nil {
		status.StorageEngine.Describe(ch)
	}
}

// GetServerStatus returns the server status info.
func GetServerStatus(session *mgo.Session) *ServerStatus {
	result := &ServerStatus{}
	err := session.DB("admin").Run(bson.D{{"serverStatus", 1}, {"recordStats", 0}}, result)
	if err != nil {
		glog.Error("Failed to get server status.")
		return nil
	}
	return result
}
