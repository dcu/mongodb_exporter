package collector

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	topTotalTimeSecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "top_total_time_seconds_total",
		Help:      "The top command provides total operation time, in seconds, and count for each database collection",
	}, []string{"type", "database", "collection"})
	topReadLockTimeSecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "top_read_lock_time_seconds_total",
		Help:      "The top command provides read lock operation time, in seconds, and count for each database collection",
	}, []string{"type", "database", "collection"})
	topWriteLockTimeSecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "top_write_lock_time_seconds_total",
		Help:      "The top command provides write lock operation time, in seconds, and count for each database collection",
	}, []string{"type", "database", "collection"})
	topQueriesTimeSecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "top_queries_time_seconds_total",
		Help:      "The top command provides queries operation time, in seconds, and count for each database collection",
	}, []string{"type", "database", "collection"})
	topGetMoreTimeSecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "top_get_more_time_seconds_total",
		Help:      "The top command provides get more operation time, in seconds, and count for each database collection",
	}, []string{"type", "database", "collection"})
	topInsertTimeSecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "top_insert_time_seconds_total",
		Help:      "The top command provides insert operation time, in seconds, and count for each database collection",
	}, []string{"type", "database", "collection"})
	topUpdateTimeSecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "top_update_time_seconds_total",
		Help:      "The top command provides update operation time, in seconds, and count for each database collection",
	}, []string{"type", "database", "collection"})
	topRemoveTimeSecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "top_remove_time_seconds_total",
		Help:      "The top command provides remove operation time, in seconds, and count for each database collection",
	}, []string{"type", "database", "collection"})
	topCommandsTimeSecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "top_commands_time_seconds_total",
		Help:      "The top command provides commands operation time, in seconds, and count for each database collection",
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

		topTotalTimeSecondsTotal.WithLabelValues("time", database, collection).Set(ConvertInSeconds(topStat.Total.Time))
		topTotalTimeSecondsTotal.WithLabelValues("count", database, collection).Set(topStat.Total.Count)

		topReadLockTimeSecondsTotal.WithLabelValues("time", database, collection).Set(ConvertInSeconds(topStat.ReadLock.Time))
		topReadLockTimeSecondsTotal.WithLabelValues("count", database, collection).Set(topStat.ReadLock.Count)

		topWriteLockTimeSecondsTotal.WithLabelValues("time", database, collection).Set(ConvertInSeconds(topStat.WriteLock.Time))
		topWriteLockTimeSecondsTotal.WithLabelValues("count", database, collection).Set(topStat.WriteLock.Count)

		topQueriesTimeSecondsTotal.WithLabelValues("time", database, collection).Set(ConvertInSeconds(topStat.Queries.Time))
		topQueriesTimeSecondsTotal.WithLabelValues("count", database, collection).Set(topStat.Queries.Count)

		topGetMoreTimeSecondsTotal.WithLabelValues("time", database, collection).Set(ConvertInSeconds(topStat.GetMore.Time))
		topGetMoreTimeSecondsTotal.WithLabelValues("count", database, collection).Set(topStat.GetMore.Count)

		topInsertTimeSecondsTotal.WithLabelValues("time", database, collection).Set(ConvertInSeconds(topStat.Insert.Time))
		topInsertTimeSecondsTotal.WithLabelValues("count", database, collection).Set(topStat.Insert.Count)

		topUpdateTimeSecondsTotal.WithLabelValues("time", database, collection).Set(ConvertInSeconds(topStat.Update.Time))
		topUpdateTimeSecondsTotal.WithLabelValues("count", database, collection).Set(topStat.Update.Count)

		topRemoveTimeSecondsTotal.WithLabelValues("time", database, collection).Set(ConvertInSeconds(topStat.Remove.Time))
		topRemoveTimeSecondsTotal.WithLabelValues("count", database, collection).Set(topStat.Remove.Count)

		topCommandsTimeSecondsTotal.WithLabelValues("time", database, collection).Set(ConvertInSeconds(topStat.Commands.Time))
		topCommandsTimeSecondsTotal.WithLabelValues("count", database, collection).Set(ConvertInSeconds(topStat.Commands.Count))
	}

	topTotalTimeSecondsTotal.Collect(ch)
	topReadLockTimeSecondsTotal.Collect(ch)
	topWriteLockTimeSecondsTotal.Collect(ch)
	topQueriesTimeSecondsTotal.Collect(ch)
	topGetMoreTimeSecondsTotal.Collect(ch)
	topInsertTimeSecondsTotal.Collect(ch)
	topUpdateTimeSecondsTotal.Collect(ch)
	topRemoveTimeSecondsTotal.Collect(ch)
	topCommandsTimeSecondsTotal.Collect(ch)
}

// Describe describes the metrics for prometheus
func (tops TopStatsMap) Describe(ch chan<- *prometheus.Desc) {
	topTotalTimeSecondsTotal.Describe(ch)
	topReadLockTimeSecondsTotal.Describe(ch)
	topWriteLockTimeSecondsTotal.Describe(ch)
	topQueriesTimeSecondsTotal.Describe(ch)
	topGetMoreTimeSecondsTotal.Describe(ch)
	topInsertTimeSecondsTotal.Describe(ch)
	topUpdateTimeSecondsTotal.Describe(ch)
	topRemoveTimeSecondsTotal.Describe(ch)
	topCommandsTimeSecondsTotal.Describe(ch)
}

// ConvertInSeconds converts microseconds in seconds (seen Prometheus best practice)
func ConvertInSeconds(microSeconds float64) float64 {
	return float64(microSeconds / 1e6)
}
