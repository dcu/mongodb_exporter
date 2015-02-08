package collector

import(
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "time"
    "github.com/prometheus/client_golang/prometheus"
)

type Collectable interface {
    Collect(exporter *MongodbCollector, ch chan<- prometheus.Metric)
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

func (status *ServerStatus) Collect(exporter *MongodbCollector, ch chan<- prometheus.Metric) {
    group := exporter.FindOrCreateGroupByName("instance")

    group.Collect(status, "Uptime", ch)
    group.Collect(status, "UptimeEstimate", ch)
    group.Collect(status, "LocalTime", ch)

    status.Asserts.Collect(exporter, ch)
}

type AssertsStats struct {
    Regular float64 `bson:"regular" type:"counter" group:"asserts"`
    Warning float64 `bson:"warning" type:"counter" group:"asserts"`
    Msg float64 `bson:"msg" type:"counter" group:"asserts"`
    User float64 `bson:"user" type:"counter" group:"asserts"`
    Rollovers float64 `bson:"rollovers" type:"counter" group:"asserts"`
}

func (asserts *AssertsStats) Collect(exporter *MongodbCollector, ch chan<- prometheus.Metric) {
    group := exporter.FindOrCreateGroupByName("asserts")

    group.Collect(asserts, "Regular", ch)
    group.Collect(asserts, "Warning", ch)
    group.Collect(asserts, "Msg", ch)
    group.Collect(asserts, "User", ch)
    group.Collect(asserts, "Rollovers", ch)
}

// Dur
type DurTiming struct {
    Dt                 float64 `bson:"dt" type:"summary" group:"durability_timing"`
    PrepLogBuffer      float64 `bson:"prepLogBuffer" type:"summary" group:"durability_timing"`
    WriteToJournal     float64 `bson:"writeToJournal" type:"summary" group:"durability_timing"`
    WriteToDataFiles   float64 `bson:"writeToDataFiles" type:"summary" group:"durability_timing"`
    RemapPrivateView   float64 `bson:"remapPrivateView" type:"summary" group:"durability_timing"`
}

type DurStats struct {
    Commits            float64 `bson:"commits" type:"gauge" group:"durability"`
    JournaledMB        float64 `bson:"journaledMB" type:"gauge" group:"durability"`
    WriteToDataFilesMB float64 `bson:"writeToDataFilesMB" type:"gauge" group:"durability"`
    Compression        float64 `bson:"compression" type:"gauge" group:"durability"`
    CommitsInWriteLock float64 `bson:"commitsInWriteLock" type:"gauge" group:"durability"`
    EarlyCommits       float64 `bson:"earlyCommits" type:"summary" group:"durability"`
    TimeMs             DurTiming `bson:"timeMs" group:"durability" type:"group"`
}

// Flush
type FlushStats struct {
    Flushes float64 `bson:"flushes" type:"counter" group:"background_flushing"`
    TotalMs float64 `bson:"total_ms" type:"counter" group:"background_flushing"`
    AverageMs float64 `bson:"average_ms" type:"gauge" group:"background_flushing"`
    LastMs float64 `bson:"last_ms" type:"gauge" group:"background_flushing"`
    LastFinished time.Time `bson:"last_finished" type:"gauge" group:"background_flushing"`
}

// Connection
type ConnectionStats struct {
    Current float64 `bson:"current" type:"gauge" group:"connections"`
    Available float64 `bson:"available" type:"gauge" group:"connections"`
    TotalCreated float64 `bson:"totalCreated" type:"counter" group:"connections"`
}

// ExtraInfo
type ExtraInfo struct {
    HeapUsageBytes float64 `bson:"heap_usage_bytes" type:"gauge" group:"extra_info"`
    PageFaults float64 `bson:"page_faults" type:"gauge" group:"extra_info"`
}

// GlobalLock
type ClientStats struct {
    Total float64 `bson:"total" type:"gauge" group:"global_lock_client"`
    Readers float64 `bson:"readers" type:"gauge" group:"global_lock_client"`
    Writers float64 `bson:"writers" type:"gauge" group:"global_lock_client"`
}

type QueueStats struct {
    Total float64 `bson:"total" type:"gauge" group:"global_lock_queue"`
    Readers float64 `bson:"readers" type:"gauge" group:"global_lock_queue"`
    Writers float64 `bson:"writers" type:"gauge" group:"global_lock_queue"`
}

type GlobalLockStats struct {
    TotalTime float64 `bson:"totalTime" type:"counter" group:"global_lock"`
    LockTime float64 `bson:"lockTime" type:"counter" group:"global_lock"`
    Ratio float64 `bson:"ratio" type:"gauge" group:"global_lock"`
    CurrentQueue *QueueStats `bson:"currentQueue" group:"global_lock" type:"group"`
    ActiveClients *ClientStats `bson:"activeClients" group:"global_lock" type:"group"`
}

//IndexCounter
type IndexCounterStats struct {
    Accesses float64 `bson:"accesses type:"counter" group:"index_counters"`
    Hits float64 `bson:"hits" type:"counter" group:"index_counters"`
    Misses float64 `bson:"misses" type:"counter" group:"index_counters"`
    Resets float64 `bson:"resets" type:"gauge" group:"index_counters"`
    MissRatio float64 `bson:"missRatio" type:"gauge" group:"index_counters"`
}

//Lock
type ReadWriteLockTimes struct {
    Read                  float64 `bson:"R" type:"counter" group:"locks"`
    Write                 float64 `bson:"W" type:"counter" group:"locks"`
    ReadLower             float64 `bson:"r" type:"counter" group:"locks"`
    WriteLower            float64 `bson:"w" type:"counter" group:"locks"`
}

type LockStats struct {
    TimeLockedMicros    ReadWriteLockTimes  `bson:"timeLockedMicros" group:"locks" type:"group"`
    TimeAcquiringMicros ReadWriteLockTimes  `bson:"timeAcquiringMicros" group:"locks" type:"group"`
}

//Network
type NetworkStats struct {
    BytesIn             float64 `bson:"bytesIn" type:"gauge" group:"network"`
    BytesOut            float64 `bson:"bytesOut" type:"gauge" group:"network"`
    NumRequests         float64 `bson:"numRequests" type:"gauge" group:"network"`
}

//Opcount and OpcountersRepl
type OpcountersStats struct {
    Insert  float64 `bson:"insert" type:"gauge" group:"op_counters"`
    Query   float64 `bson:"query" type:"gauge" group:"op_counters"`
    Update  float64 `bson:"update" type:"gauge" group:"op_counters"`
    Delete  float64 `bson:"delete" type:"gauge" group:"op_counters"`
    GetMore float64 `bson:"getmore" type:"gauge" group:"op_counters"`
    Command float64 `bson:"command" type:"gauge" group:"op_counters"`
}

//Mem
type MemStats struct {
    Bits                float64 `bson:"bits" type:"counter" group:"memory"`
    Resident            float64 `bson:"resident" type:"counter" group:"memory"`
    Virtual             float64 `bson:"virtual" type:"gauge" group:"memory"`
    Mapped              float64 `bson:"mapped" type:"gauge" group:"memory"`
    MappedWithJournal   float64 `bson:"mappedWithJournal" type:"counter" group:"memory"`
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

