package collector

import(
    "testing"
	"github.com/dcu/mongodb_exporter/shared"
)

func Test_ExtraInfoCollectData(t *testing.T) {
	stats := &ExtraInfo{
	}

	groupName := "extra_info"
	stats.Collect(groupName, nil)

	if shared.Groups[groupName] == nil {
		t.Error("Group not created")
	}
}

