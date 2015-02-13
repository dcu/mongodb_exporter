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

	asserts.Collect("asserts", nil)

	if shared.Groups["asserts"] == nil {
		t.Error("Asserts group was not loaded.")
	}
}

