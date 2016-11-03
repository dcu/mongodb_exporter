package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	locksTimeLockedGlobalMicrosecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "locks_time_locked_global_microseconds_total",
		Help:      "amount of time in microseconds that any database has held the global lock",
	}, []string{"type", "database"})
	locksTimeLockedLocalMicrosecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "locks_time_locked_local_microseconds_total",
		Help:      "amount of time in microseconds that any database has held the local lock",
	}, []string{"type", "database"})
	locksTimeAcquiringGlobalMicrosecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "locks_time_acquiring_global_microseconds_total",
		Help:      "amount of time in microseconds that any database has spent waiting for the global lock",
	}, []string{"type", "database"})
)

// LockStatsMap is a map of lock stats
type LockStatsMap map[string]LockStats

// ReadWriteLockTimes information about the lock
type ReadWriteLockTimes struct {
	Read       float64 `bson:"R"`
	Write      float64 `bson:"W"`
	ReadLower  float64 `bson:"r"`
	WriteLower float64 `bson:"w"`
}

// LockStats lock stats
type LockStats struct {
	TimeLockedMicros    ReadWriteLockTimes `bson:"timeLockedMicros"`
	TimeAcquiringMicros ReadWriteLockTimes `bson:"timeAcquiringMicros"`
}

// Export exports the data to prometheus.
func (locks LockStatsMap) Export(ch chan<- prometheus.Metric) {
	for key, locks := range locks {
		if key == "." {
			key = "dot"
		}

		locksTimeLockedGlobalMicrosecondsTotal.WithLabelValues("read", key).Add(locks.TimeLockedMicros.Read)
		locksTimeLockedGlobalMicrosecondsTotal.WithLabelValues("write", key).Add(locks.TimeLockedMicros.Write)

		locksTimeLockedLocalMicrosecondsTotal.WithLabelValues("read", key).Add(locks.TimeLockedMicros.ReadLower)
		locksTimeLockedLocalMicrosecondsTotal.WithLabelValues("write", key).Add(locks.TimeLockedMicros.WriteLower)

		locksTimeAcquiringGlobalMicrosecondsTotal.WithLabelValues("read", key).Add(locks.TimeAcquiringMicros.ReadLower)
		locksTimeAcquiringGlobalMicrosecondsTotal.WithLabelValues("write", key).Add(locks.TimeAcquiringMicros.WriteLower)
	}

	locksTimeLockedGlobalMicrosecondsTotal.Collect(ch)
	locksTimeLockedLocalMicrosecondsTotal.Collect(ch)
	locksTimeAcquiringGlobalMicrosecondsTotal.Collect(ch)
}

// Describe describes the metrics for prometheus
func (locks LockStatsMap) Describe(ch chan<- *prometheus.Desc) {
	locksTimeLockedGlobalMicrosecondsTotal.Describe(ch)
	locksTimeLockedLocalMicrosecondsTotal.Describe(ch)
	locksTimeAcquiringGlobalMicrosecondsTotal.Describe(ch)
}
