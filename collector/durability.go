package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
)

// Dur
type DurTiming struct {
	Dt               float64 `bson:"dt"`
	PrepLogBuffer    float64 `bson:"prepLogBuffer"`
	WriteToJournal   float64 `bson:"writeToJournal"`
	WriteToDataFiles float64 `bson:"writeToDataFiles"`
	RemapPrivateView float64 `bson:"remapPrivateView"`
}

func (durTiming *DurTiming) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("dt", durTiming.Dt)
	group.Export("prep_log_buffer", durTiming.PrepLogBuffer)
	group.Export("write_to_journal", durTiming.WriteToJournal)
	group.Export("write_to_data_files", durTiming.WriteToDataFiles)
	group.Export("remap_private_view", durTiming.RemapPrivateView)
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

func (durStats *DurStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName + "_commits")
	group.Export("written", durStats.Commits)
	group.Export("in_write_lock", durStats.CommitsInWriteLock)

	group = shared.FindOrCreateGroup(groupName)
	group.Export("journaled_megabytes", durStats.JournaledMB)
	group.Export("write_to_data_files_megabytes", durStats.WriteToDataFilesMB)
	group.Export("compression", durStats.Compression)
	group.Export("early_commits", durStats.EarlyCommits)

	durStats.TimeMs.Export(groupName + "_time_milliseconds")

}
