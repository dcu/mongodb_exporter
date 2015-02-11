package collector

import(
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "time"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/dcu/mongodb_exporter/shared"
)

type ServerStatus struct {
    Uptime             float64                `bson:"uptime"`
    UptimeEstimate     float64                `bson:"uptimeEstimate"`
    LocalTime          time.Time              `bson:"localTime"`

    Asserts            *AssertsStats          `bson:"asserts"`

    Dur                *DurStats              `bson:"dur"`

    BackgroundFlushing *FlushStats            `bson:"backgroundFlushing"`

    Connections        *ConnectionStats       `bson:"connections"`

    ExtraInfo          *ExtraInfo             `bson:"extra_info"`

    GlobalLock         *GlobalLockStats       `bson:"globalLock"`

    IndexCounter       *IndexCounterStats     `bson:"indexCounters"`

    Locks              LockStatsMap   `bson:"locks,omitempty"`

    Network            *NetworkStats          `bson:"network"`

    Opcounters         *OpcountersStats       `bson:"opcounters"`
    OpcountersRepl     *OpcountersStats       `bson:"opcountersRepl"`
    Mem                *MemStats              `bson:"mem"`
    Metrics            *MetricsStats          `bson:"metrics"`
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
    status.Metrics.Collect("metrics", ch)
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

