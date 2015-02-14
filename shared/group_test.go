package shared

import(
	"testing"
	"github.com/prometheus/client_golang/prometheus"
)

func Test_FindOrCreateGroup(t *testing.T) {
	group := FindOrCreateGroup("new_group")

	if group == nil {
		t.Error("Group couldn't be created.")
	}
}

func Test_Collect(t *testing.T) {
	LoadGroupsDesc()

	group := FindOrCreateGroup("asserts")
	group.Collect("regular", 10.0, nil)
}

func Test_GetGauge(t *testing.T) {
	chDesc := make(chan *prometheus.Desc)
	chCollect := make(chan prometheus.Metric)

	LoadGroupsDesc()
	group := FindOrCreateGroup("background_flushing")
	go group.Collect("average_ms", 1.0, chCollect)
	go group.Describe(chDesc)
}

func Test_GetCounter(t *testing.T) {
	chDesc := make(chan *prometheus.Desc)
	chCollect := make(chan prometheus.Metric)

	group := FindOrCreateGroup("asserts")
	go group.Collect("regular", 1.0, chCollect)
	go group.Describe(chDesc)
}

func Test_GetSummary(t *testing.T) {
	chDesc := make(chan *prometheus.Desc)
	chCollect := make(chan prometheus.Metric)

	group := FindOrCreateGroup("durability")
	go group.Collect("early_commits", 1.0, chCollect)
	go group.Describe(chDesc)
}

