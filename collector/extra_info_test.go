package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"testing"
)

func Test_ExtraInfoCollectData(t *testing.T) {
	stats := &ExtraInfo{}

	groupName := "extra_info"
	stats.Export(groupName)

	if shared.Groups[groupName] == nil {
		t.Error("Group not created")
	}
}
