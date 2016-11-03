package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	opCountersTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "op_counters_total",
		Help:      "The opcounters data structure provides an overview of database operations by type and makes it possible to analyze the load on the database in more granular manner. These numbers will grow over time and in response to database use. Analyze these values over time to track database utilization",
	}, []string{"type"})
	opCountersReplTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "op_counters_repl_total",
		Help:      "The opcountersRepl data structure, similar to the opcounters data structure, provides an overview of database replication operations by type and makes it possible to analyze the load on the replica in more granular manner. These values only appear when the current host has replication enabled",
	}, []string{"type"})
)

// OpcountersStats opcounters stats
type OpcountersStats struct {
	Insert  float64 `bson:"insert"`
	Query   float64 `bson:"query"`
	Update  float64 `bson:"update"`
	Delete  float64 `bson:"delete"`
	GetMore float64 `bson:"getmore"`
	Command float64 `bson:"command"`
}

// Export exports the data to prometheus.
func (opCounters *OpcountersStats) Export(ch chan<- prometheus.Metric) {
	opCountersTotal.WithLabelValues("insert").Add(opCounters.Insert)
	opCountersTotal.WithLabelValues("query").Add(opCounters.Query)
	opCountersTotal.WithLabelValues("update").Add(opCounters.Update)
	opCountersTotal.WithLabelValues("delete").Add(opCounters.Delete)
	opCountersTotal.WithLabelValues("getmore").Add(opCounters.GetMore)
	opCountersTotal.WithLabelValues("command").Add(opCounters.Command)

	opCountersTotal.Collect(ch)
}

// Describe describes the metrics for prometheus
func (opCounters *OpcountersStats) Describe(ch chan<- *prometheus.Desc) {
	opCountersTotal.Describe(ch)
}

// OpcountersReplStats opcounters stats
type OpcountersReplStats struct {
	Insert  float64 `bson:"insert"`
	Query   float64 `bson:"query"`
	Update  float64 `bson:"update"`
	Delete  float64 `bson:"delete"`
	GetMore float64 `bson:"getmore"`
	Command float64 `bson:"command"`
}

// Export exports the data to prometheus.
func (opCounters *OpcountersReplStats) Export(ch chan<- prometheus.Metric) {
	opCountersReplTotal.WithLabelValues("insert").Add(opCounters.Insert)
	opCountersReplTotal.WithLabelValues("query").Add(opCounters.Query)
	opCountersReplTotal.WithLabelValues("update").Add(opCounters.Update)
	opCountersReplTotal.WithLabelValues("delete").Add(opCounters.Delete)
	opCountersReplTotal.WithLabelValues("getmore").Add(opCounters.GetMore)
	opCountersReplTotal.WithLabelValues("command").Add(opCounters.Command)

	opCountersReplTotal.Collect(ch)
}

// Describe describes the metrics for prometheus
func (opCounters *OpcountersReplStats) Describe(ch chan<- *prometheus.Desc) {
	opCountersReplTotal.Describe(ch)
}
