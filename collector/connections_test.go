package collector

import(
    "testing"
	"github.com/dcu/mongodb_exporter/shared"
)

func Test_ConnectionsCollectData(t *testing.T) {
	stats := &ConnectionStats{
		Current: 1,
		Available: 2,
		TotalCreated: 3,
	}

	groupName := "connections"
	stats.Export(groupName)

	if shared.Groups[groupName] == nil {
		t.Error("Group not created")
	}
}

