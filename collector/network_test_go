package collector

import(
    "testing"
    "github.com/dcu/mongodb_exporter/shared"
)

func Test_NetworkCollectData(t *testing.T) {
	stats := &NetworkStats{
	}

	groupName := "network"
	stats.Collect(groupName, nil)

	if shared.Groups[groupName] == nil {
		t.Error("Group not created")
	}
}

