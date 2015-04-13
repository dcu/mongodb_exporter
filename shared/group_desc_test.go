package shared

import (
	"testing"
)

func Test_LoadGroupDesc(t *testing.T) {
	LoadGroupsDesc()

	if len(GroupsDesc) == 0 {
		t.Error("Groups were not loaded.")
	}
}

func Test_ParseEnabledGroups(t *testing.T) {
	ParseEnabledGroups("a, b,  c")
	if !EnabledGroups["a"] {
		t.Error("a was not loaded.")
	}
	if !EnabledGroups["b"] {
		t.Error("b was not loaded.")
	}
	if !EnabledGroups["c"] {
		t.Error("c was not loaded.")
	}
}

func Test_GroupField(t *testing.T) {
	LoadGroupsDesc()

	field := GroupField("instance", "uptime_seconds")

	if field.Type != "counter" {
		t.Error("field was not loaded.")
	}
}
