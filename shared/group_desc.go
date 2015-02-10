package shared

import(
    "gopkg.in/yaml.v2"
)


type FieldDesc struct {
    Type string
    Help string
}

type GroupFieldsMap map[string]*FieldDesc
type GroupDescMap map[string]GroupFieldsMap

var(
    GroupsDesc = make(GroupDescMap)
)

func LoadGroupsDesc() {
    dat, errx := Asset("groups.yml")
    if(errx != nil){panic(errx)}
    yaml.Unmarshal(dat, GroupsDesc)
}

func GroupFields(groupName string) GroupFieldsMap {
    fields := GroupsDesc[groupName]
    if fields == nil {
        panic("Couldn't find group:"+groupName)
    }

    return fields
}

func GroupField(groupName string, fieldName string) *FieldDesc {
    field := GroupFields(groupName)[fieldName]

    if field == nil {
        panic("Couldn't find field: "+fieldName+" in: "+groupName)
    }

    return field
}

