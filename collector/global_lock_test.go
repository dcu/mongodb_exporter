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
	stats.Export(groupName)

	if shared.Groups[groupName] == nil {
		t.Error("Group not created")
	}
}

