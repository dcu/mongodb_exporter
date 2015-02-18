package collector

import(
    "testing"
	"github.com/dcu/mongodb_exporter/shared"
)

func Test_AssertsCollectData(t *testing.T) {
	asserts := &AssertsStats{
		Regular: 1,
		Warning: 2,
		Msg: 3,
		User: 4,
		Rollovers: 5,
	}

	asserts.Export("asserts")

	if shared.Groups["asserts_total"] == nil {
		t.Error("Asserts group was not loaded.")
	}
}

