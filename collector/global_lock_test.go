package collector

import(
    "testing"
	"github.com/dcu/mongodb_exporter/shared"
)

func Test_GlobalLockCollectData(t *testing.T) {
	stats := &GlobalLockStats{
		CurrentQueue: &QueueStats{},
		ActiveClients: &ClientStats{},
	}

	groupName := "global_lock"
	stats.Collect(groupName, nil)

	if shared.Groups[groupName] == nil {
		t.Error("Group not created")
	}
}

