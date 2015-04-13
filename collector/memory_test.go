package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"testing"
)

func Test_MemoryCollectData(t *testing.T) {
	stats := &MemStats{}

	groupName := "memory"
	stats.Export(groupName)

	if shared.Groups[groupName] == nil {
		t.Error("Group not created")
	}
}
