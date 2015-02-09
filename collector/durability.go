package collector

import(
    "github.com/prometheus/client_golang/prometheus"
)


// Dur
type DurTiming struct {
    Dt                 float64 `bson:"dt" type:"summary"`
    PrepLogBuffer      float64 `bson:"prepLogBuffer" type:"summary"`
    WriteToJournal     float64 `bson:"writeToJournal" type:"summary"`
    WriteToDataFiles   float64 `bson:"writeToDataFiles" type:"summary"`
    RemapPrivateView   float64 `bson:"remapPrivateView" type:"summary"`
}
func (durTiming *DurTiming) Collect(groupName string, exporter *MongodbCollector, ch chan<- prometheus.Metric) {
    group := exporter.FindOrCreateGroupByName(groupName)
    group.Collect(durTiming, "Dt", ch)
    group.Collect(durTiming, "PrepLogBuffer", ch)
    group.Collect(durTiming, "WriteToJournal", ch)
    group.Collect(durTiming, "WriteToDataFiles", ch)
    group.Collect(durTiming, "RemapPrivateView", ch)
}

type DurStats struct {
    Commits            float64 `bson:"commits" type:"gauge"`
    JournaledMB        float64 `bson:"journaledMB" type:"gauge"`
    WriteToDataFilesMB float64 `bson:"writeToDataFilesMB" type:"gauge"`
    Compression        float64 `bson:"compression" type:"gauge"`
    CommitsInWriteLock float64 `bson:"commitsInWriteLock" type:"gauge"`
    EarlyCommits       float64 `bson:"earlyCommits" type:"summary"`
    TimeMs             DurTiming `bson:"timeMs"`
}

func (durStats *DurStats) Collect(groupName string, exporter *MongodbCollector, ch chan<- prometheus.Metric) {
    group := exporter.FindOrCreateGroupByName(groupName)
    group.Collect(durStats, "Commits", ch)
    group.Collect(durStats, "JournaledMB", ch)
    group.Collect(durStats, "WriteToDataFilesMB", ch)
    group.Collect(durStats, "Compression", ch)
    group.Collect(durStats, "CommitsInWriteLock", ch)
    group.Collect(durStats, "EarlyCommits", ch)

    durStats.TimeMs.Collect(groupName+"_time_ms", exporter, ch)
}


