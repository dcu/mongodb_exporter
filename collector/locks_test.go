package collector

import(
    "testing"
	"github.com/dcu/mongodb_exporter/shared"
)

func Test_LocksCollectData(t *testing.T) {
	stats := &LockStatsMap{
		".": LockStats{
			TimeLockedMicros: ReadWriteLockTimes{},
			TimeAcquiringMicros: ReadWriteLockTimes{},
		},
	}

	groupName := "locks"
	stats.Collect(groupName, nil)

	if shared.Groups["dot_locks_time_locked_global"] == nil {
		t.Error("Group not created")
	}
	if shared.Groups["dot_locks_time_acquiring_global"] == nil {
		t.Error("Group not created")
	}
}

