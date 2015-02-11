package collector

import(
    "github.com/prometheus/client_golang/prometheus"
    "github.com/dcu/mongodb_exporter/shared"
)

type LockStatsMap map[string]LockStats

//Lock
type ReadWriteLockTimes struct {
    Read                  float64 `bson:"R"`
    Write                 float64 `bson:"W"`
    ReadLower             float64 `bson:"r"`
    WriteLower            float64 `bson:"w"`
}

type LockStats struct {
    TimeLockedMicros    ReadWriteLockTimes  `bson:"timeLockedMicros"`
    TimeAcquiringMicros ReadWriteLockTimes  `bson:"timeAcquiringMicros"`
}

func (locks LockStatsMap) Collect(groupName string, ch chan<-prometheus.Metric) {
    for key, locks := range locks {
        if key == "." {
            key = "dot"
        }

        timeLockedGroup := shared.FindOrCreateGroup(key+"_locks_time_locked")
        timeLockedGroup.DescName = "locks_time_locked_micros"

        timeLockedGroup.Collect("global_r", locks.TimeLockedMicros.Read, ch)
        timeLockedGroup.Collect("global_w", locks.TimeLockedMicros.Write, ch)
        timeLockedGroup.Collect("local_r", locks.TimeLockedMicros.ReadLower, ch)
        timeLockedGroup.Collect("local_w", locks.TimeLockedMicros.WriteLower, ch)

        timeAcquiringGroup := shared.FindOrCreateGroup(key+"_locks_time_acquiring")
        timeAcquiringGroup.DescName = "locks_time_acquiring_micros"
        timeAcquiringGroup.Collect("global_r", locks.TimeLockedMicros.Read, ch)
        timeAcquiringGroup.Collect("global_w", locks.TimeLockedMicros.Write, ch)
    }
}

