package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"testing"
)

func Test_DurabilityCollectData(t *testing.T) {
	stats := &DurStats{
		TimeMs: DurTiming{},
	}

	groupName := "durability"
	stats.Export(groupName)

	if shared.Groups[groupName] == nil {
		t.Error("Group not created")
	}
}
