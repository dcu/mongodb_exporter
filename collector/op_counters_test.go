package collector

import(
    "testing"
	"github.com/dcu/mongodb_exporter/shared"
)

func Test_OpCountersCollectData(t *testing.T) {
	stats := &OpcountersStats{
	}

	groupName := "op_counters"
	stats.Export(groupName)

	if shared.Groups[groupName+"_total"] == nil {
		t.Error("Group not created")
	}
}

