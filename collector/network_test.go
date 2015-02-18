package collector

import(
    "testing"
    "github.com/dcu/mongodb_exporter/shared"
)

func Test_NetworkCollectData(t *testing.T) {
	stats := &NetworkStats{
	}

	groupName := "network"
	stats.Export(groupName)

	if shared.Groups[groupName+"_bytes_total"] == nil {
		t.Error("Group not created")
	}
	if shared.Groups[groupName+"_metrics"] == nil {
		t.Error("Group not created")
	}
}

