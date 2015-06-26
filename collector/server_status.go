package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
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
func (status *ServerStatus) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)

	group.Export("uptime_seconds", status.Uptime)
	group.Export("uptime_estimate_seconds", status.Uptime)
	group.Export("local_time", float64(status.LocalTime.Unix()))

	if status.Asserts != nil {
		exportData(status.Asserts, "asserts")
	}

	if status.Dur != nil {
		exportData(status.Dur, "durability")
	}

	if status.BackgroundFlushing != nil {
		exportData(status.BackgroundFlushing, "background_flushing")
	}

	if status.Connections != nil {
		exportData(status.Connections, "connections")
	}

	if status.ExtraInfo != nil {
		exportData(status.ExtraInfo, "extra_info")
	}

	if status.GlobalLock != nil {
		exportData(status.GlobalLock, "global_lock")
	}

	if status.IndexCounter != nil {
		exportData(status.IndexCounter, "index_counters")
	}

	if status.Network != nil {
		exportData(status.Network, "network")
	}

	if status.Opcounters != nil {
		exportData(status.Opcounters, "op_counters")
	}

	if status.OpcountersRepl != nil {
		exportData(status.OpcountersRepl, "op_counters_repl")
	}

	if status.Mem != nil {
		exportData(status.Mem, "memory")
	}

	if status.Locks != nil {
		exportData(status.Locks, "locks")
	}

	if status.Metrics != nil {
		exportData(status.Metrics, "metrics")
	}

	if status.Cursors != nil {
		exportData(status.Cursors, "cursors")
	}
}

func exportData(exportable shared.Exportable, groupName string) {
	if !shared.EnabledGroups[groupName] {
		// disabled group
		glog.Infof("Group is not enabled: %s", groupName)
		return
	}

	exportable.Export(groupName)
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
