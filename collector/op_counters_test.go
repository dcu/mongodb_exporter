package collector

import (
	"testing"
)

func Test_OpCountersCollectData(t *testing.T) {
	stats := &OpcountersStats{}

	stats.Export()
}
