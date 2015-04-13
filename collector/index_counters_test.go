package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"testing"
)

func Test_IndexCountersCollectData(t *testing.T) {
	stats := &IndexCounterStats{}

	groupName := "index_counters"
	stats.Export(groupName)

	if shared.Groups[groupName] == nil {
		t.Error("Group not created")
	}
}
