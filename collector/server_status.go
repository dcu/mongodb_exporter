package collector

import(
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "time"
    "github.com/prometheus/client_golang/prometheus"
)

type Collectable interface {
    Collect(groupName string, exporter *MongodbCollector, ch chan<- prometheus.Metric)
}

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

    Locks              map[string]LockStats   `bson:"locks,omitempty" group:"locks" type:"group"`

    Network            *NetworkStats          `bson:"network" group:"network" type:"group"`

    Opcounters         *OpcountersStats       `bson:"opcounters" group:"op_counters" type:"group"`
    OpcountersRepl     *OpcountersStats       `bson:"opcountersRepl" group:"op_counters_repl" type:"group"`
    Mem                *MemStats              `bson:"mem" group:"memory" type:"group"`
}

func (status *ServerStatus) Collect(groupName string, exporter *MongodbCollector, ch chan<- prometheus.Metric) {
    group := exporter.FindOrCreateGroupByName(groupName)

    group.Collect(status, "Uptime", ch)
    group.Collect(status, "UptimeEstimate", ch)
    group.Collect(status, "LocalTime", ch)

    status.Asserts.Collect("asserts", exporter, ch)
    status.Dur.Collect("durability", exporter, ch)
    status.BackgroundFlushing.Collect("background_flushing", exporter, ch)
    status.Connections.Collect("connections", exporter, ch)
    status.ExtraInfo.Collect("extra_info", exporter, ch)
    status.GlobalLock.Collect("global_lock", exporter, ch)
    status.IndexCounter.Collect("index_counters", exporter, ch)
    status.Network.Collect("network", exporter, ch)
    status.Opcounters.Collect("op_counters", exporter, ch)
    status.OpcountersRepl.Collect("op_counters_repl", exporter, ch)
    status.Mem.Collect("memory", exporter, ch)
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

