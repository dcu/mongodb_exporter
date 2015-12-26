package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	durabilityjournaledMegabytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "durability",
		Name:      "journaled_megabytes",
		Help:      "The journaledMB provides the amount of data in megabytes (MB) written to journal during the last journal group commit interval",
	})
	durabilitywriteToDataFilesMegabytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "durability",
		Name:      "write_to_data_files_megabytes",
		Help:      "The writeToDataFilesMB provides the amount of data in megabytes (MB) written from journal to the data files during the last journal group commit interval",
	})
	durabilitycompression = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "durability",
		Name:      "compression",
		Help:      "The compression represents the compression ratio of the data written to the journal: ( journaled_size_of_data / uncompressed_size_of_data )",
	})
	durabilityearlyCommits = prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: Namespace,
		Subsystem: "durability",
		Name:      "early_commits",
		Help:      "The earlyCommits value reflects the number of times MongoDB requested a commit before the scheduled journal group commit interval. Use this value to ensure that your journal group commit interval is not too long for your deployment",
	})
)

var (
	durabilityCommits = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "durability_commits",
		Help:      "Durability commits",
	}, []string{})
)

var (
	durabilityTimeMilliseconds = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: Namespace,
		Name:      "durability_time_milliseconds",
		Help:      "Summary of times spent during the journaling process.",
	}, []string{})
)

// DurTiming is the information about durability returned from the server.
type DurTiming struct {
	Dt               float64 `bson:"dt"`
	PrepLogBuffer    float64 `bson:"prepLogBuffer"`
	WriteToJournal   float64 `bson:"writeToJournal"`
	WriteToDataFiles float64 `bson:"writeToDataFiles"`
	RemapPrivateView float64 `bson:"remapPrivateView"`
}

// Export exports the data for the prometheus server.
func (durTiming *DurTiming) Export() {
	durabilityTimeMilliseconds.WithLabelValues("dt").Observe(durTiming.Dt)
	durabilityTimeMilliseconds.WithLabelValues("prep_log_buffer").Observe(durTiming.PrepLogBuffer)
	durabilityTimeMilliseconds.WithLabelValues("write_to_journal").Observe(durTiming.WriteToJournal)
	durabilityTimeMilliseconds.WithLabelValues("write_to_data_files").Observe(durTiming.WriteToDataFiles)
	durabilityTimeMilliseconds.WithLabelValues("remap_private_view").Observe(durTiming.RemapPrivateView)
}

// DurStats are the stats related to durability.
type DurStats struct {
	Commits            float64   `bson:"commits"`
	JournaledMB        float64   `bson:"journaledMB"`
	WriteToDataFilesMB float64   `bson:"writeToDataFilesMB"`
	Compression        float64   `bson:"compression"`
	CommitsInWriteLock float64   `bson:"commitsInWriteLock"`
	EarlyCommits       float64   `bson:"earlyCommits"`
	TimeMs             DurTiming `bson:"timeMs"`
}

// Export export the durability stats for the prometheus server.
func (durStats *DurStats) Export() {
	durabilityCommits.WithLabelValues("written").Add(durStats.Commits)
	durabilityCommits.WithLabelValues("in_write_lock").Add(durStats.CommitsInWriteLock)

	durabilityjournaledMegabytes.Add(durStats.JournaledMB)
	durabilitywriteToDataFilesMegabytes.Add(durStats.WriteToDataFilesMB)
	durabilitycompression.Add(durStats.Compression)
	durabilityearlyCommits.Observe(durStats.EarlyCommits)

	durStats.TimeMs.Export()
}
