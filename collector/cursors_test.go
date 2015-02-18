package collector

import(
    "testing"
	"github.com/dcu/mongodb_exporter/shared"
)

func Test_CursorsCollectData(t *testing.T) {
	cursors := &Cursors{
	}

	cursors.Export("cursors")
	if shared.Groups["cursors"] == nil {
		t.Error("Cursors group was not loaded.")
	}
}

