package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"github.com/prometheus/client_golang/prometheus"
)

// Dur
type DurTiming struct {
	Dt               float64 `bson:"dt"`
	PrepLogBuffer    float64 `bson:"prepLogBuffer"`
	WriteToJournal   float64 `bson:"writeToJournal"`
	WriteToDataFiles float64 `bson:"writeToDataFiles"`
	RemapPrivateView float64 `bson:"remapPrivateView"`
}

func (durTiming *DurTiming) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)
	group.Collect("dt", durTiming.Dt, ch)
	group.Collect("prep_log_buffer", durTiming.PrepLogBuffer, ch)
	group.Collect("write_to_journal", durTiming.WriteToJournal, ch)
	group.Collect("write_to_data_files", durTiming.WriteToDataFiles, ch)
	group.Collect("remap_private_view", durTiming.RemapPrivateView, ch)
}

type DurStats struct {
	Commits            float64   `bson:"commits"`
	JournaledMB        float64   `bson:"journaledMB"`
	WriteToDataFilesMB float64   `bson:"writeToDataFilesMB"`
	Compression        float64   `bson:"compression"`
	CommitsInWriteLock float64   `bson:"commitsInWriteLock"`
	EarlyCommits       float64   `bson:"earlyCommits"`
	TimeMs             DurTiming `bson:"timeMs"`
}

func (durStats *DurStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)
	group.Collect("commits", durStats.Commits, ch)
	group.Collect("journaled_mb", durStats.JournaledMB, ch)
	group.Collect("write_to_data_files_mb", durStats.WriteToDataFilesMB, ch)
	group.Collect("compression", durStats.Compression, ch)
	group.Collect("commits_in_write_lock", durStats.CommitsInWriteLock, ch)
	group.Collect("early_commits", durStats.EarlyCommits, ch)

	durStats.TimeMs.Collect(groupName+"_time_ms", ch)
}
