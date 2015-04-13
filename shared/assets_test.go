package shared

import (
	"testing"
)

func Test_GroupsYml_Data(t *testing.T) {
	data, err := Asset("groups.yml")

	if err != nil {
		panic(err)
	}

	if data == nil || len(data) == 0 {
		t.Error("groups.yml was not found")
	}
}

func Test_GroupsYml_Info(t *testing.T) {
	info, err := AssetInfo("groups.yml")
	if err != nil {
		panic(err)
	}

	if info.Size() < 24305 {
		t.Error("Loaded asset is too small.")
	}
}
