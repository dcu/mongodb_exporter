package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/golang/glog"
)

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

	Cursors        *Cursors         `bson:"cursors"`
}

func (status *ServerStatus) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)

	group.Collect("uptime", status.Uptime, ch)
	group.Collect("uptime_estimate", status.Uptime, ch)
	group.Collect("local_time", float64(status.LocalTime.Unix()), ch)

	collectData(status.Asserts, "asserts", ch)
	collectData(status.Dur, "durability", ch)
	collectData(status.BackgroundFlushing, "background_flushing", ch)
	collectData(status.Connections, "connections", ch)
	collectData(status.ExtraInfo, "extra_info", ch)
	collectData(status.GlobalLock, "global_lock", ch)
	collectData(status.IndexCounter, "index_counters", ch)
	collectData(status.Network, "network", ch)
	collectData(status.Opcounters, "op_counters", ch)
	collectData(status.OpcountersRepl, "op_counters_repl", ch)
	collectData(status.Mem, "memory", ch)
	collectData(status.Locks, "locks", ch)
	collectData(status.Metrics, "metrics", ch)
	collectData(status.Cursors, "cursors", ch)
}

func collectData(collectable shared.Collectable, groupName string, ch chan<- prometheus.Metric) {
	if !shared.EnabledGroups[groupName] {
		// disabled group
		return
	}

	collectable.Collect(groupName, ch)
}

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

