package collector_mongos

import (
	"testing"
)

func Test_IndexCountersCollectData(t *testing.T) {
	stats := &IndexCounterStats{}

	stats.Export()
}
