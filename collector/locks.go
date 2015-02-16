package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
)

type LockStatsMap map[string]LockStats

//Lock
type ReadWriteLockTimes struct {
	Read       float64 `bson:"R"`
	Write      float64 `bson:"W"`
	ReadLower  float64 `bson:"r"`
	WriteLower float64 `bson:"w"`
}

type LockStats struct {
	TimeLockedMicros    ReadWriteLockTimes `bson:"timeLockedMicros"`
	TimeAcquiringMicros ReadWriteLockTimes `bson:"timeAcquiringMicros"`
}

func (locks LockStatsMap) Export(groupName string) {
	for key, locks := range locks {
		if key == "." {
			key = "dot"
		}

		timeLockedGroup := shared.FindOrCreateGroup(key + "_locks_time_locked_microseconds_global")
		timeLockedGroup.DescName = "locks_time_locked_global_microseconds_total"
		timeLockedGroup.Export("read", locks.TimeLockedMicros.Read)
		timeLockedGroup.Export("write", locks.TimeLockedMicros.Write)

		timeLockedGroup = shared.FindOrCreateGroup(key + "_locks_time_locked_microseconds_local")
		timeLockedGroup.DescName = "locks_time_locked_local_microseconds_total"
		timeLockedGroup.Export("read", locks.TimeLockedMicros.ReadLower)
		timeLockedGroup.Export("write", locks.TimeLockedMicros.WriteLower)

		timeAcquiringGroup := shared.FindOrCreateGroup(key + "_locks_time_acquiring_microseconds_global")
		timeAcquiringGroup.DescName = "locks_time_acquiring_global_microseconds_total"
		timeAcquiringGroup.Export("read", locks.TimeLockedMicros.Read)
		timeAcquiringGroup.Export("write", locks.TimeLockedMicros.Write)
	}
}
