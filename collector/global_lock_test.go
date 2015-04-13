package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"testing"
)

func Test_GlobalLockCollectData(t *testing.T) {
	stats := &GlobalLockStats{
		CurrentQueue:  &QueueStats{},
		ActiveClients: &ClientStats{},
	}

	groupName := "global_lock"
	stats.Export(groupName)

	if shared.Groups[groupName] == nil {
		t.Error("Group not created")
	}
}
