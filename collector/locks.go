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

		timeLockedGroup := shared.FindOrCreateGroup("locks_time_locked_microseconds_global")
		timeLockedGroup.DescName = "locks_time_locked_global_microseconds_total"
		timeLockedGroup.ExportWithLabels([]string{"read", key}, locks.TimeLockedMicros.Read)
		timeLockedGroup.ExportWithLabels([]string{"write", key}, locks.TimeLockedMicros.Write)

		timeLockedGroup = shared.FindOrCreateGroup("locks_time_locked_microseconds_local")
		timeLockedGroup.DescName = "locks_time_locked_local_microseconds_total"
		timeLockedGroup.ExportWithLabels([]string{"read", key}, locks.TimeLockedMicros.ReadLower)
		timeLockedGroup.ExportWithLabels([]string{"write", key}, locks.TimeLockedMicros.WriteLower)

		timeAcquiringGroup := shared.FindOrCreateGroup("locks_time_acquiring_microseconds_global")
		timeAcquiringGroup.DescName = "locks_time_acquiring_global_microseconds_total"
		timeAcquiringGroup.ExportWithLabels([]string{"read", key}, locks.TimeLockedMicros.Read)
		timeAcquiringGroup.ExportWithLabels([]string{"write", key}, locks.TimeLockedMicros.Write)
	}
}
