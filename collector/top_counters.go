package collector

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	topTotalTimeMicrosecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "top_total_time_microseconds_total",
		Help:      "The top command provides total operation time, in microseconds, and count for each database collection",
	}, []string{"type", "database", "collection"})
	topReadLockTimeMicrosecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "top_read_lock_time_microseconds_total",
		Help:      "The top command provides read lock operation time, in microseconds, and count for each database collection",
	}, []string{"type", "database", "collection"})
	topWriteLockTimeMicrosecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "top_write_lock_time_microseconds_total",
		Help:      "The top command provides write lock operation time, in microseconds, and count for each database collection",
	}, []string{"type", "database", "collection"})
	topQueriesTimeMicrosecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "top_queries_time_microseconds_total",
		Help:      "The top command provides queries operation time, in microseconds, and count for each database collection",
	}, []string{"type", "database", "collection"})
	topGetMoreTimeMicrosecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "top_get_more_time_microseconds_total",
		Help:      "The top command provides get more operation time, in microseconds, and count for each database collection",
	}, []string{"type", "database", "collection"})
	topInsertTimeMicrosecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "top_insert_time_microseconds_total",
		Help:      "The top command provides insert operation time, in microseconds, and count for each database collection",
	}, []string{"type", "database", "collection"})
	topUpdateTimeMicrosecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "top_update_time_microseconds_total",
		Help:      "The top command provides update operation time, in microseconds, and count for each database collection",
	}, []string{"type", "database", "collection"})
	topRemoveTimeMicrosecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "top_remove_time_microseconds_total",
		Help:      "The top command provides remove operation time, in microseconds, and count for each database collection",
	}, []string{"type", "database", "collection"})
	topCommandsTimeMicrosecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "top_commands_time_microseconds_total",
		Help:      "The top command provides commands operation time, in microseconds, and count for each database collection",
	}, []string{"type", "database", "collection"})
)

// TopStatsMap is a map of top stats
type TopStatsMap map[string]TopStats

// TopcountersStats topcounters stats
type TopcounterStats struct {
	Time  float64 `bson:"time"`
	Count float64 `bson:"count"`
}

// TopCollectionStats top collection stats
type TopStats struct {
	Total     TopcounterStats `bson:"total"`
	ReadLock  TopcounterStats `bson:"readLock"`
	WriteLock TopcounterStats `bson:"writeLock"`
	Queries   TopcounterStats `bson:"queries"`
	GetMore   TopcounterStats `bson:"getmore"`
	Insert    TopcounterStats `bson:"insert"`
	Update    TopcounterStats `bson:"update"`
	Remove    TopcounterStats `bson:"remove"`
	Commands  TopcounterStats `bson:"commands"`
}

// Export exports the data to prometheus.
func (topStats TopStatsMap) Export(ch chan<- prometheus.Metric) {

	for collectionNamespace, topStat := range topStats {

		namespace := strings.Split(collectionNamespace, ".")
		database := namespace[0]
		collection := strings.Join(namespace[1:], ".")

		topTotalTimeMicrosecondsTotal.WithLabelValues("time", database, collection).Set(topStat.Total.Time)
		topTotalTimeMicrosecondsTotal.WithLabelValues("count", database, collection).Set(topStat.Total.Count)

		topReadLockTimeMicrosecondsTotal.WithLabelValues("time", database, collection).Set(topStat.ReadLock.Time)
		topReadLockTimeMicrosecondsTotal.WithLabelValues("count", database, collection).Set(topStat.ReadLock.Count)

		topWriteLockTimeMicrosecondsTotal.WithLabelValues("time", database, collection).Set(topStat.WriteLock.Time)
		topWriteLockTimeMicrosecondsTotal.WithLabelValues("count", database, collection).Set(topStat.WriteLock.Count)

		topQueriesTimeMicrosecondsTotal.WithLabelValues("time", database, collection).Set(topStat.Queries.Time)
		topQueriesTimeMicrosecondsTotal.WithLabelValues("count", database, collection).Set(topStat.Queries.Count)

		topGetMoreTimeMicrosecondsTotal.WithLabelValues("time", database, collection).Set(topStat.GetMore.Time)
		topGetMoreTimeMicrosecondsTotal.WithLabelValues("count", database, collection).Set(topStat.GetMore.Count)

		topInsertTimeMicrosecondsTotal.WithLabelValues("time", database, collection).Set(topStat.Insert.Time)
		topInsertTimeMicrosecondsTotal.WithLabelValues("count", database, collection).Set(topStat.Insert.Count)

		topUpdateTimeMicrosecondsTotal.WithLabelValues("time", database, collection).Set(topStat.Update.Time)
		topUpdateTimeMicrosecondsTotal.WithLabelValues("count", database, collection).Set(topStat.Update.Count)

		topRemoveTimeMicrosecondsTotal.WithLabelValues("time", database, collection).Set(topStat.Remove.Time)
		topRemoveTimeMicrosecondsTotal.WithLabelValues("count", database, collection).Set(topStat.Remove.Count)

		topCommandsTimeMicrosecondsTotal.WithLabelValues("time", database, collection).Set(topStat.Commands.Time)
		topCommandsTimeMicrosecondsTotal.WithLabelValues("count", database, collection).Set(topStat.Commands.Count)
	}

	topTotalTimeMicrosecondsTotal.Collect(ch)
	topReadLockTimeMicrosecondsTotal.Collect(ch)
	topWriteLockTimeMicrosecondsTotal.Collect(ch)
	topQueriesTimeMicrosecondsTotal.Collect(ch)
	topGetMoreTimeMicrosecondsTotal.Collect(ch)
	topInsertTimeMicrosecondsTotal.Collect(ch)
	topUpdateTimeMicrosecondsTotal.Collect(ch)
	topRemoveTimeMicrosecondsTotal.Collect(ch)
	topCommandsTimeMicrosecondsTotal.Collect(ch)
}

// Describe describes the metrics for prometheus
func (tops TopStatsMap) Describe(ch chan<- *prometheus.Desc) {
	topTotalTimeMicrosecondsTotal.Describe(ch)
	topReadLockTimeMicrosecondsTotal.Describe(ch)
	topWriteLockTimeMicrosecondsTotal.Describe(ch)
	topQueriesTimeMicrosecondsTotal.Describe(ch)
	topGetMoreTimeMicrosecondsTotal.Describe(ch)
	topInsertTimeMicrosecondsTotal.Describe(ch)
	topUpdateTimeMicrosecondsTotal.Describe(ch)
	topRemoveTimeMicrosecondsTotal.Describe(ch)
	topCommandsTimeMicrosecondsTotal.Describe(ch)
}
