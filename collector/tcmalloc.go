package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	tcmallocGeneral = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "tcmalloc_generic_heap",
		Help:      "High-level summary metricsInternal metrics from tcmalloc",
	}, []string{"type"})
	tcmallocPageheapBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "tcmalloc_pageheap_bytes",
		Help:      "Sizes for tcpmalloc pageheaps",
	}, []string{"type"})
	tcmallocPageheapCounts = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "tcmalloc_pageheap_count",
		Help:      "Sizes for tcpmalloc pageheaps",
	}, []string{"type"})

	tcmallocCacheBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "tcmalloc_cache_bytes",
		Help:      "Sizes for tcpmalloc caches in bytes",
	}, []string{"cache", "type"})

	tcmallocAggressiveDecommit = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "tcmalloc_aggressive_memory_decommit",
		Help:      "Whether aggressive_memory_decommit is on",
	})

	tcmallocFreeBytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "tcmalloc_free_bytes",
		Help:      "Total free bytes of tcmalloc",
	})
)

// TCMallocStats tracks the mem stats metrics.
type TCMallocStats struct {
	Generic GenericTCMAllocStats  `bson:"generic"`
	Details DetailedTCMallocStats `bson:"tcmalloc"`
}

type GenericTCMAllocStats struct {
	CurrentAllocatedBytes float64 `bson:"current_allocated_bytes"`
	HeapSize              float64 `bson:"heap_size"`
}

type DetailedTCMallocStats struct {
	PageheapFreeBytes          float64 `bson:"pageheap_free_bytes"`
	PageheapUnmappedBytes      float64 `bson:"pageheap_unmapped_bytes"`
	PageheapComittedBytes      float64 `bson:"pageheap_committed_bytes"`
	PageheapScavengeCount      float64 `bson:"pageheap_scavenge_count"`
	PageheapCommitCount        float64 `bson:"pageheap_commit_count"`
	PageheapTotalCommitBytes   float64 `bson:"pageheap_total_commit_bytes"`
	PageheapDecommitCount      float64 `bson:"pageheap_decommit_count"`
	PageheapTotalDecommitBytes float64 `bson:"pageheap_total_decommit_bytes"`
	PageheapReserveCount       float64 `bson:"pageheap_reserve_count"`
	PageheapTotalReserveBytes  float64 `bson:"pageheap_total_reserve_bytes"`

	MaxTotalThreadCacheBytes     float64 `bson:"max_total_thread_cache_bytes"`
	CurrentTotalThreadCacheBytes float64 `bson:"current_total_thread_cache_bytes"`
	CentralCacheFreeBytes        float64 `bson:"central_cache_free_bytes"`
	TransferCacheFreeBytes       float64 `bson:"transfer_cache_free_bytes"`
	ThreadCacheFreeBytes         float64 `bson:"thread_cache_free_bytes"`

	TotalFreeBytes           float64 `bson:"total_free_bytes"`
	AggressiveMemoryDecommit float64 `bson:"aggressive_memory_decommit"`
}

// Export exports the data to prometheus.
func (m *TCMallocStats) Export(ch chan<- prometheus.Metric) {
	// Generic metrics
	tcmallocGeneral.WithLabelValues("allocated").Set(m.Generic.CurrentAllocatedBytes)
	tcmallocGeneral.WithLabelValues("total").Set(m.Generic.HeapSize)
	tcmallocGeneral.Collect(ch)

	// Pageheap
	tcmallocPageheapBytes.WithLabelValues("free").Set(m.Details.PageheapFreeBytes)
	tcmallocPageheapBytes.WithLabelValues("unmapped").Set(m.Details.PageheapUnmappedBytes)
	tcmallocPageheapBytes.WithLabelValues("comitted").Set(m.Details.PageheapComittedBytes)
	tcmallocPageheapBytes.WithLabelValues("total_commit").Set(m.Details.PageheapTotalCommitBytes)
	tcmallocPageheapBytes.WithLabelValues("total_decommit").Set(m.Details.PageheapTotalDecommitBytes)
	tcmallocPageheapBytes.WithLabelValues("total_reserve").Set(m.Details.PageheapTotalReserveBytes)
	tcmallocPageheapBytes.Collect(ch)

	tcmallocPageheapCounts.WithLabelValues("scavenge").Set(m.Details.PageheapScavengeCount)
	tcmallocPageheapCounts.WithLabelValues("commit").Set(m.Details.PageheapCommitCount)
	tcmallocPageheapCounts.WithLabelValues("decommit").Set(m.Details.PageheapDecommitCount)
	tcmallocPageheapCounts.WithLabelValues("reserve").Set(m.Details.PageheapReserveCount)
	tcmallocPageheapCounts.Collect(ch)

	tcmallocCacheBytes.WithLabelValues("thread_cache", "max_total").Set(m.Details.MaxTotalThreadCacheBytes)
	tcmallocCacheBytes.WithLabelValues("thread_cache", "current_total").Set(m.Details.CurrentTotalThreadCacheBytes)
	tcmallocCacheBytes.WithLabelValues("central_cache", "free").Set(m.Details.CentralCacheFreeBytes)
	tcmallocCacheBytes.WithLabelValues("transfer_cache", "free").Set(m.Details.TransferCacheFreeBytes)
	tcmallocCacheBytes.WithLabelValues("thread_cache", "free").Set(m.Details.ThreadCacheFreeBytes)
	tcmallocCacheBytes.Collect(ch)

	tcmallocAggressiveDecommit.Set(m.Details.AggressiveMemoryDecommit)
	tcmallocAggressiveDecommit.Collect(ch)

	tcmallocFreeBytes.Set(m.Details.TotalFreeBytes)
	tcmallocFreeBytes.Collect(ch)

}

// Describe describes the metrics for prometheus
func (m *TCMallocStats) Describe(ch chan<- *prometheus.Desc) {
	tcmallocGeneral.Describe(ch)
	tcmallocPageheapBytes.Describe(ch)
	tcmallocPageheapCounts.Describe(ch)
	tcmallocCacheBytes.Describe(ch)
	tcmallocAggressiveDecommit.Describe(ch)
	tcmallocFreeBytes.Describe(ch)
}
