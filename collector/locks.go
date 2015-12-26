package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	locksTimeLockedGlobalMicrosecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "locks_time_locked_global_microseconds_total",
		Help:      "amount of time in microseconds that any database has held the global lock",
	}, []string{})
)
var (
	locksTimeLockedLocalMicrosecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "locks_time_locked_local_microseconds_total",
		Help:      "amount of time in microseconds that any database has held the local lock",
	}, []string{})
)
var (
	locksTimeAcquiringGlobalMicrosecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "locks_time_acquiring_global_microseconds_total",
		Help:      "amount of time in microseconds that any database has spent waiting for the global lock",
	}, []string{})
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
func (locks LockStatsMap) Export() {
	for key, locks := range locks {
		if key == "." {
			key = "dot"
		}

		locksTimeLockedGlobalMicrosecondsTotal.WithLabelValues("read", key).Set(locks.TimeLockedMicros.Read)
		locksTimeLockedGlobalMicrosecondsTotal.WithLabelValues("write", key).Set(locks.TimeLockedMicros.Write)

		locksTimeLockedLocalMicrosecondsTotal.WithLabelValues("read", key).Set(locks.TimeLockedMicros.ReadLower)
		locksTimeLockedLocalMicrosecondsTotal.WithLabelValues("write", key).Set(locks.TimeLockedMicros.WriteLower)

		locksTimeAcquiringGlobalMicrosecondsTotal.WithLabelValues("read", key).Set(locks.TimeAcquiringMicros.ReadLower)
		locksTimeAcquiringGlobalMicrosecondsTotal.WithLabelValues("write", key).Set(locks.TimeAcquiringMicros.WriteLower)
	}
}
