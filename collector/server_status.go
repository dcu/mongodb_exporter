package collector

import(
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "time"
)

type ServerStatus struct {
    Uptime             float64                `bson:"uptime" type:"counter"`
    UptimeEstimate     float64                `bson:"uptimeEstimate" type:"counter"`
    LocalTime          time.Time              `bson:"localTime" type:"counter"`

    Asserts            *AssertsStats          `bson:"asserts"`

    Dur                *DurStats              `bson:"dur"`

    BackgroundFlushing *FlushStats            `bson:"backgroundFlushing"`

    Connections        *ConnectionStats       `bson:"connections"`

    ExtraInfo          *ExtraInfo             `bson:"extra_info"`

    GlobalLock         *GlobalLockStats       `bson:"globalLock"`

    IndexCounter       *IndexCounterStats     `bson:"indexCounters"`

    Locks              map[string]LockStats   `bson:"locks,omitempty"`

    Network            *NetworkStats          `bson:"network"`

    Opcounters         *OpcountStats          `bson:"opcounters"`
    OpcountersRepl     *OpcountStats          `bson:"opcountersRepl"`
    Mem                *MemStats              `bson:"mem"`
}

type AssertsStats struct {
    Regular            float64                `bson:"regular" type:"counter"`
    Warning            float64                `bson:"warning type:"counter""`
    Msg                float64                `bson:"msg type:"counter""`
    User               float64                `bson:"user type:"counter""`
    Rollovers          float64                `bson:"rollovers" type:"counter"`
}

// Dur
type DurTiming struct {
    Dt                 float64                `bson:"dt" type:"summary"`
    PrepLogBuffer      float64                `bson:"prepLogBuffer" type:"summary"`
    WriteToJournal     float64                `bson:"writeToJournal" type:"summary"`
    WriteToDataFiles   float64                `bson:"writeToDataFiles" type:"summary"`
    RemapPrivateView   float64                `bson:"remapPrivateView" type:"summary"`
}

type DurStats struct {
    Commits            float64              `bson:"commits" type:"gauge"`
    JournaledMB        float64              `bson:"journaledMB" type:"gauge"`
    WriteToDataFilesMB float64              `bson:"writeToDataFilesMB" type:"gauge"`
    Compression        float64              `bson:"compression" type:"gauge"`
    CommitsInWriteLock float64              `bson:"commitsInWriteLock" type:"gauge"`
    EarlyCommits       float64              `bson:"earlyCommits" type:"summary"`
    TimeMs           DurTiming              `bson:"timeMs"`
}

// Flush
type FlushStats struct {
    Flushes            float64              `bson:"flushes" type:"counter"`
    TotalMs            float64              `bson:"total_ms" type:"counter"`
    AverageMs          float64              `bson:"average_ms" type:"gauge"`
    LastMs             float64              `bson:"last_ms" type:"gauge"`
    LastFinished     time.Time              `bson:"last_finished" type:"gauge"`
}

// Connection
type ConnectionStats struct {
    Current            float64              `bson:"current" type:"gauge"`
    Available          float64              `bson:"available" type:"gauge"`
    TotalCreated       float64              `bson:"totalCreated" type:"counter"`
}

// ExtraInfo
type ExtraInfo struct {
    HeapUsageBytes     float64              `bson:"heap_usage_bytes" type:"gauge"`
    PageFaults         float64              `bson:"page_faults" type:"gauge"`
}

// GlobalLock
type ClientStats struct {
    Total                float64            `bson:"total" type:"gauge"`
    Readers              float64            `bson:"readers" type:"gauge"`
    Writers              float64            `bson:"writers" type:"gauge"`
}

type QueueStats struct {
    Total                float64            `bson:"total" type:"gauge"`
    Readers              float64            `bson:"readers" type:"gauge"`
    Writers              float64            `bson:"writers" type:"gauge"`
}

type GlobalLockStats struct {
    TotalTime             float64           `bson:"totalTime" type:"counter"`
    LockTime              float64           `bson:"lockTime" type:"counter"`
    Ratio                 float64           `bson:"ratio" type:"gauge"`
    CurrentQueue          *QueueStats       `bson:"currentQueue"`
    ActiveClients         *ClientStats      `bson:"activeClients"`
}

//IndexCounter
type IndexCounterStats struct {
    Accesses              float64           `bson:"accesses type:"counter""`
    Hits                  float64           `bson:"hits" type:"counter"`
    Misses                float64           `bson:"misses" type:"counter"`
    Resets                float64           `bson:"resets" type:"gauge"`
    MissRatio             float64           `bson:"missRatio" type:"gauge"`
}

//Lock
type ReadWriteLockTimes struct {
    Read                  float64           `bson:"R" type:"counter"`
    Write                 float64           `bson:"W" type:"counter"`
    ReadLower             float64           `bson:"r" type:"counter"`
    WriteLower            float64           `bson:"w" type:"counter"`
}

type LockStats struct {
    TimeLockedMicros    ReadWriteLockTimes  `bson:"timeLockedMicros"`
    TimeAcquiringMicros ReadWriteLockTimes  `bson:"timeAcquiringMicros"`
}

//Network
type NetworkStats struct {
    BytesIn             float64             `bson:"bytesIn" type:"gauge"`
    BytesOut            float64             `bson:"bytesOut" type:"gauge"`
    NumRequests         float64             `bson:"numRequests" type:"gauge"`
}

//Opcount and OpcountersRepl
type OpcountStats struct {
    Insert              float64             `bson:"insert" type:"gauge"`
    Query               float64             `bson:"query" type:"gauge"`
    Update              float64             `bson:"update" type:"gauge"`
    Delete              float64             `bson:"delete" type:"gauge"`
    GetMore             float64             `bson:"getmore" type:"gauge"`
    Command             float64             `bson:"command" type:"gauge"`
}

//Mem
type MemStats struct {
    Bits                float64             `bson:"bits" type:"counter"`
    Resident            float64             `bson:"resident" type:"counter"`
    Virtual             float64             `bson:"virtual" type:"gauge"`
    Supported           interface{}         `bson:"supported"`
    Mapped              float64             `bson:"mapped" type:"gauge"`
    MappedWithJournal   float64             `bson:"mappedWithJournal" type:"counter"`
}

func NewServerStatus() *ServerStatus {
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

    return result
}

