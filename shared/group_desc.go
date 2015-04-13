package shared

import (
	"gopkg.in/yaml.v2"
	"strings"
)

type FieldDesc struct {
	Type   string
	Labels []string
	Help   string
}

type GroupFieldsMap map[string]*FieldDesc
type GroupDescMap map[string]GroupFieldsMap

var (
	GroupsDesc    = make(GroupDescMap)
	EnabledGroups = make(map[string]bool)
)

func LoadGroupsDesc() {
	dat, err := Asset("groups.yml")
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal(dat, GroupsDesc)
}

func ParseEnabledGroups(enabledGroupsFlag string) {
	for _, name := range strings.Split(enabledGroupsFlag, ",") {
		name = strings.TrimSpace(name)
		EnabledGroups[name] = true
	}
}

func GroupFields(groupName string) GroupFieldsMap {
	fields := GroupsDesc[groupName]
	if fields == nil {
		panic("Couldn't find group:" + groupName)
	}

	return fields
}

func GroupField(groupName string, fieldName string) *FieldDesc {
	field := GroupFields(groupName)[fieldName]

	if field == nil {
		panic("Couldn't find field: " + fieldName + " in: " + groupName)
	}

	return field
}
