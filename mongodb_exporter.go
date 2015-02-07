package main

import(
    "net/http"
    "github.com/prometheus/client_golang/prometheus"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "time"
    "fmt"
)

type ServerStatus struct {
  Uptime             int64                  `bson:"uptime"`
  UptimeEstimate     int64                  `bson:"uptimeEstimate"`
  LocalTime          time.Time              `bson:"localTime"`

  Asserts            map[string]int64       `bson:"asserts"`

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
  Dt               int64 `bson:"dt"`
  PrepLogBuffer    int64 `bson:"prepLogBuffer"`
  WriteToJournal   int64 `bson:"writeToJournal"`
  WriteToDataFiles int64 `bson:"writeToDataFiles"`
  RemapPrivateView int64 `bson:"remapPrivateView"`
}

type DurStats struct {
  Commits            int64 `bson:"commits"`
  JournaledMB        int64 `bson:"journaledMB"`
  WriteToDataFilesMB int64 `bson:"writeToDataFilesMB"`
  Compression        int64 `bson:"compression"`
  CommitsInWriteLock int64 `bson:"commitsInWriteLock"`
  EarlyCommits       int64 `bson:"earlyCommits"`
  TimeMs             DurTiming `bson:"timeMs"`
}

// Flush
type FlushStats struct {
  Flushes      int64     `bson:"flushes"`
  TotalMs      int64     `bson:"total_ms"`
  AverageMs    float64   `bson:"average_ms"`
  LastMs       int64     `bson:"last_ms"`
  LastFinished time.Time `bson:"last_finished"`
}

// Connection
type ConnectionStats struct {
  Current      int64 `bson:"current"`
  Available    int64 `bson:"available"`
  TotalCreated int64 `bson:"totalCreated"`
}

// ExtraInfo
type ExtraInfo struct {
  PageFaults *int64 `bson:"page_faults"`
}

// GlobalLock
type QueueStats struct {
  Total   int64 `bson:"total"`
  Readers int64 `bson:"readers"`
  Writers int64 `bson:"writers"`
}

type ClientStats struct {
  Total   int64 `bson:"total"`
  Readers int64 `bson:"readers"`
  Writers int64 `bson:"writers"`
}

type GlobalLockStats struct {
  TotalTime     int64        `bson:"totalTime"`
  LockTime      int64        `bson:"lockTime"`
  CurrentQueue  *QueueStats  `bson:"currentQueue"`
  ActiveClients *ClientStats `bson:"activeClients"`
}

//IndexCounter
type IndexCounterStats struct {
  Accesses  int64 `bson:"accesses"`
  Hits      int64 `bson:"hits"`
  Misses    int64 `bson:"misses"`
  Resets    int64 `bson:"resets"`
  MissRatio int64 `bson:"missRatio"`
}

//Lock
type ReadWriteLockTimes struct {
  Read       int64 `bson:"R"`
  Write      int64 `bson:"W"`
  ReadLower  int64 `bson:"r"`
  WriteLower int64 `bson:"w"`
}

type LockStats struct {
  TimeLockedMicros    ReadWriteLockTimes `bson:"timeLockedMicros"`
  TimeAcquiringMicros ReadWriteLockTimes `bson:"timeAcquiringMicros"`
}

//Network
type NetworkStats struct {
  BytesIn     int64 `bson:"bytesIn"`
  BytesOut    int64 `bson:"bytesOut"`
  NumRequests int64 `bson:"numRequests"`
}

//Opcount and OpcountersRepl
type OpcountStats struct {
  Insert  int64 `bson:"insert"`
  Query   int64 `bson:"query"`
  Update  int64 `bson:"update"`
  Delete  int64 `bson:"delete"`
  GetMore int64 `bson:"getmore"`
  Command int64 `bson:"command"`
}

//Mem
type MemStats struct {
  Bits              int64       `bson:"bits"`
  Resident          int64       `bson:"resident"`
  Virtual           int64       `bson:"virtual"`
  Supported         interface{} `bson:"supported"`
  Mapped            int64       `bson:"mapped"`
  MappedWithJournal int64       `bson:"mappedWithJournal"`
}

func main() {
    result := &ServerStatus{}

    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }
    defer session.Close()

    session.SetMode(mgo.Eventual, true)
    session.SetSocketTimeout(0)
    defer session.Close()

    err = session.DB("admin").Run(bson.D{{"serverStatus", 1}, {"recordStats", 0}}, result)
    fmt.Println("serverStatus:", result)

    http.Handle("/metrics", prometheus.Handler())
    http.ListenAndServe(":9001", nil)
}
