package collector

import(
    "testing"
	"github.com/dcu/mongodb_exporter/shared"
)

func Test_BackgroundFlushingCollectData(t *testing.T) {
	stats := &FlushStats{
		Flushes: 1,
		TotalMs: 2,
		AverageMs: 3,
		LastMs: 4,
	}

	groupName := "background_flushing"
	stats.Export(groupName)

	if shared.Groups[groupName] == nil {
		t.Error("Group not created")
	}
}

