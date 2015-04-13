package shared

import (
	"github.com/prometheus/client_golang/prometheus"
	"testing"
)

func Test_FindOrCreateGroup(t *testing.T) {
	group := FindOrCreateGroup("new_group")

	if group == nil {
		t.Error("Group couldn't be created.")
	}
}

func Test_Collect(t *testing.T) {
	LoadGroupsDesc()

	group := FindOrCreateGroup("asserts_total")
	group.Export("regular", 10.0)
}

func Test_GetGauge(t *testing.T) {
	chDesc := make(chan *prometheus.Desc)
	chCollect := make(chan prometheus.Metric)

	LoadGroupsDesc()
	group := FindOrCreateGroup("background_flushing")
	group.Export("average_milliseconds", 1.0)
	go group.Collect(chCollect)
	go group.Describe(chDesc)
}

func Test_GetCounter(t *testing.T) {
	chDesc := make(chan *prometheus.Desc)
	chCollect := make(chan prometheus.Metric)

	group := FindOrCreateGroup("asserts_total")
	group.Export("regular", 1.0)
	go group.Collect(chCollect)
	go group.Describe(chDesc)
}

func Test_GetSummary(t *testing.T) {
	chDesc := make(chan *prometheus.Desc)
	chCollect := make(chan prometheus.Metric)

	group := FindOrCreateGroup("durability")
	group.Export("early_commits", 1.0)
	go group.Collect(chCollect)
	go group.Describe(chDesc)
}
