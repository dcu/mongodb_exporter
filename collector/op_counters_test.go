package collector

import(
    "testing"
	"github.com/dcu/mongodb_exporter/shared"
)

func Test_OpCountersCollectData(t *testing.T) {
	stats := &OpcountersStats{
	}

	groupName := "op_counters"
	stats.Collect(groupName, nil)

	if shared.Groups[groupName] == nil {
		t.Error("Group not created")
	}
}

