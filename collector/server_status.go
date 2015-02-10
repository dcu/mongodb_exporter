package collector

import(
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "time"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/dcu/mongodb_exporter/shared"
)

type ServerStatus struct {
    Uptime             float64                `bson:"uptime" group:"instance" type:"counter"`
    UptimeEstimate     float64                `bson:"uptimeEstimate" group:"instance" type:"counter"`
    LocalTime          time.Time              `bson:"localTime" group:"instance" type:"counter"`

    Asserts            *AssertsStats          `bson:"asserts" group:"asserts" type:"group"`

    Dur                *DurStats              `bson:"dur" group:"durability" type:"group"`

    BackgroundFlushing *FlushStats            `bson:"backgroundFlushing" group:"background_flushing" type:"group"`

    Connections        *ConnectionStats       `bson:"connections" group:"connections" type:"group"`

    ExtraInfo          *ExtraInfo             `bson:"extra_info" group:"extra_info" type:"group"`

    GlobalLock         *GlobalLockStats       `bson:"globalLock" group:"global_lock" type:"group"`

    IndexCounter       *IndexCounterStats     `bson:"indexCounters" group:"index_counters" type:"group"`

    Locks              LockStatsMap   `bson:"locks,omitempty" group:"locks" type:"group"`

    Network            *NetworkStats          `bson:"network" group:"network" type:"group"`

    Opcounters         *OpcountersStats       `bson:"opcounters" group:"op_counters" type:"group"`
    OpcountersRepl     *OpcountersStats       `bson:"opcountersRepl" group:"op_counters_repl" type:"group"`
    Mem                *MemStats              `bson:"mem" group:"memory" type:"group"`
}

func (status *ServerStatus) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)

    group.Collect("uptime", status.Uptime, ch)
    group.Collect("uptime_estimate", status.Uptime, ch)
    group.Collect("local_time", float64(status.LocalTime.Unix()), ch)

    status.Asserts.Collect("asserts", ch)
    status.Dur.Collect("durability", ch)
    status.BackgroundFlushing.Collect("background_flushing", ch)
    status.Connections.Collect("connections", ch)
    status.ExtraInfo.Collect("extra_info", ch)
    status.GlobalLock.Collect("global_lock", ch)
    status.IndexCounter.Collect("index_counters", ch)
    status.Network.Collect("network", ch)
    status.Opcounters.Collect("op_counters", ch)
    status.OpcountersRepl.Collect("op_counters_repl", ch)
    status.Mem.Collect("memory", ch)
    status.Locks.Collect("locks", ch)
}

func GetServerStatus() *ServerStatus {
    result := &ServerStatus{}
    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }

    session.SetMode(mgo.Eventual, true)
    session.SetSocketTimeout(0)
    defer func() {
        println("Closing connection to database.")
        session.Close()
    }()

    err = session.DB("admin").Run(bson.D{{"serverStatus", 1}, {"recordStats", 0}}, result)
    if err != nil {
        println("Failed to get server status.")
    }

    return result
}

