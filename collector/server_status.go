package collector

import(
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "time"
    "fmt"
)

type ServerStatus struct {
  Uptime             float64                  `bson:"uptime"`
  UptimeEstimate     float64                  `bson:"uptimeEstimate"`
  LocalTime          time.Time              `bson:"localTime"`

  Asserts            map[string]float64       `bson:"asserts"`

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

// Dur
type DurTiming struct {
  Dt               float64 `bson:"dt"`
  PrepLogBuffer    float64 `bson:"prepLogBuffer"`
  WriteToJournal   float64 `bson:"writeToJournal"`
  WriteToDataFiles float64 `bson:"writeToDataFiles"`
  RemapPrivateView float64 `bson:"remapPrivateView"`
}

type DurStats struct {
  Commits            float64 `bson:"commits"`
  JournaledMB        float64 `bson:"journaledMB"`
  WriteToDataFilesMB float64 `bson:"writeToDataFilesMB"`
  Compression        float64 `bson:"compression"`
  CommitsInWriteLock float64 `bson:"commitsInWriteLock"`
  EarlyCommits       float64 `bson:"earlyCommits"`
  TimeMs             DurTiming `bson:"timeMs"`
}

// Flush
type FlushStats struct {
  Flushes      float64     `bson:"flushes"`
  TotalMs      float64     `bson:"total_ms"`
  AverageMs    float64   `bson:"average_ms"`
  LastMs       float64     `bson:"last_ms"`
  LastFinished time.Time `bson:"last_finished"`
}

// Connection
type ConnectionStats struct {
  Current      float64 `bson:"current"`
  Available    float64 `bson:"available"`
  TotalCreated float64 `bson:"totalCreated"`
}

// ExtraInfo
type ExtraInfo struct {
  PageFaults *float64 `bson:"page_faults"`
}

// GlobalLock
type QueueStats struct {
  Total   float64 `bson:"total"`
  Readers float64 `bson:"readers"`
  Writers float64 `bson:"writers"`
}

type ClientStats struct {
  Total   float64 `bson:"total"`
  Readers float64 `bson:"readers"`
  Writers float64 `bson:"writers"`
}

type GlobalLockStats struct {
  TotalTime     float64        `bson:"totalTime"`
  LockTime      float64        `bson:"lockTime"`
  CurrentQueue  *QueueStats  `bson:"currentQueue"`
  ActiveClients *ClientStats `bson:"activeClients"`
}

//IndexCounter
type IndexCounterStats struct {
  Accesses  float64 `bson:"accesses"`
  Hits      float64 `bson:"hits"`
  Misses    float64 `bson:"misses"`
  Resets    float64 `bson:"resets"`
  MissRatio float64 `bson:"missRatio"`
}

//Lock
type ReadWriteLockTimes struct {
  Read       float64 `bson:"R"`
  Write      float64 `bson:"W"`
  ReadLower  float64 `bson:"r"`
  WriteLower float64 `bson:"w"`
}

type LockStats struct {
  TimeLockedMicros    ReadWriteLockTimes `bson:"timeLockedMicros"`
  TimeAcquiringMicros ReadWriteLockTimes `bson:"timeAcquiringMicros"`
}

//Network
type NetworkStats struct {
  BytesIn     float64 `bson:"bytesIn"`
  BytesOut    float64 `bson:"bytesOut"`
  NumRequests float64 `bson:"numRequests"`
}

//Opcount and OpcountersRepl
type OpcountStats struct {
  Insert  float64 `bson:"insert"`
  Query   float64 `bson:"query"`
  Update  float64 `bson:"update"`
  Delete  float64 `bson:"delete"`
  GetMore float64 `bson:"getmore"`
  Command float64 `bson:"command"`
}

//Mem
type MemStats struct {
  Bits              float64       `bson:"bits"`
  Resident          float64       `bson:"resident"`
  Virtual           float64       `bson:"virtual"`
  Supported         interface{} `bson:"supported"`
  Mapped            float64       `bson:"mapped"`
  MappedWithJournal float64       `bson:"mappedWithJournal"`
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
    fmt.Println("serverStatus:", result.Connections.Current)

    return result
}

