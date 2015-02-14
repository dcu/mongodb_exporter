package shared

import(
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

	group := FindOrCreateGroup("asserts")
	group.Collect("regular", 10.0, nil)
}

func Test_GetGauge(t *testing.T) {
	LoadGroupsDesc()
	group := FindOrCreateGroup("background_flushing")
	group.Collect("average_ms", 1.0, nil)
}

func Test_GetCounter(t *testing.T) {
	group := FindOrCreateGroup("asserts")
	group.Collect("regular", 1.0, nil)
}
func Test_GetSummary(t *testing.T) {
	group := FindOrCreateGroup("durability")
	group.Collect("early_commits", 1.0, nil)
}

