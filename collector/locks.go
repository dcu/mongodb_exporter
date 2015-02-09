package collector

import(
    "github.com/prometheus/client_golang/prometheus"
)

type LockStatsMap map[string]LockStats

//Lock
type ReadWriteLockTimes struct {
    Read                  float64 `bson:"R" type:"counter"`
    Write                 float64 `bson:"W" type:"counter"`
    ReadLower             float64 `bson:"r" type:"counter"`
    WriteLower            float64 `bson:"w" type:"counter"`
}

type LockStats struct {
    TimeLockedMicros    ReadWriteLockTimes  `bson:"timeLockedMicros"`
    TimeAcquiringMicros ReadWriteLockTimes  `bson:"timeAcquiringMicros"`
}

func (locks LockStatsMap) Collect(groupName string, exporter *MongodbCollector, ch chan<-prometheus.Metric) {
    for key, locks := range locks {
        if key == "." {
            key = "dot"
        }

        timeLockedGroup := exporter.FindOrCreateGroupByName(key+"_locks_time_locked")
        timeLockedGroup.Collect(locks.TimeLockedMicros, "Read", ch)
        timeLockedGroup.Collect(locks.TimeLockedMicros, "Write", ch)
        timeLockedGroup.Collect(locks.TimeLockedMicros, "ReadLower", ch)
        timeLockedGroup.Collect(locks.TimeLockedMicros, "WriteLower", ch)

        timeAcquiringGroup := exporter.FindOrCreateGroupByName(key+"_locks_time_acquiring")
        timeAcquiringGroup.Collect(locks.TimeLockedMicros, "Read", ch)
        timeAcquiringGroup.Collect(locks.TimeLockedMicros, "Write", ch)
        timeAcquiringGroup.Collect(locks.TimeLockedMicros, "ReadLower", ch)
        timeAcquiringGroup.Collect(locks.TimeLockedMicros, "WriteLower", ch)
    }
}

