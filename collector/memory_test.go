package collector

import(
    "testing"
	"github.com/dcu/mongodb_exporter/shared"
)

func Test_MemoryCollectData(t *testing.T) {
	stats := &MemStats{
	}

	groupName := "memory"
	stats.Export(groupName)

	if shared.Groups[groupName] == nil {
		t.Error("Group not created")
	}
}

